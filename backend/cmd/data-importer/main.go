package main

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/atdevten/peace/internal/infrastructure/config"
	"github.com/atdevten/peace/internal/infrastructure/database"
	"github.com/atdevten/peace/internal/infrastructure/database/postgres/repository"
)

type QuoteJob struct {
	Content string
	Author  string
	Row     int
}

type WorkerStats struct {
	Processed *atomic.Uint32
	Errors    *atomic.Uint32
	Skipped   *atomic.Uint32
}

type FileProcessor interface {
	ProcessFile(filePath string) error
}

type CSVProcessor struct {
	quoteRepo repositories.QuoteRepository
	dryRun    bool
	workers   int
	batchSize int
	stats     *WorkerStats
}

func (p *CSVProcessor) ProcessFile(filePath string) error {
	// Open CSV file
	csvFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file %s: %w", filePath, err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1

	ctx := context.Background()

	log.Printf("Starting to process CSV file: %s", filePath)
	log.Printf("Using %d workers with batch size %d", p.workers, p.batchSize)
	if p.dryRun {
		log.Println("DRY RUN mode enabled - no data will be inserted")
	}

	// Create job channel and start workers
	jobChan := make(chan QuoteJob, p.batchSize)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go worker(ctx, i+1, jobChan, p.quoteRepo, p.dryRun, p.stats, &wg)
	}

	// Start progress reporter
	done := make(chan bool)
	go progressReporter(p.stats, done)

	// Skip header row
	if _, err := csvReader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Read CSV and send jobs to workers
	rowCount := 0
	startTime := time.Now()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV row %d: %v", rowCount+1, err)
			p.stats.Errors.Add(1)
			continue
		}

		rowCount++

		// Parse CSV record (quote, author, category)
		if len(record) < 2 {
			log.Printf("Skipping invalid record at row %d (insufficient columns)", rowCount)
			p.stats.Skipped.Add(1)
			continue
		}

		quoteContent := strings.TrimSpace(record[0])
		author := strings.TrimSpace(record[1])

		// Skip empty quotes
		if quoteContent == "" {
			log.Printf("Skipping empty quote at row %d", rowCount)
			p.stats.Skipped.Add(1)
			continue
		}

		// Set default author if empty
		if author == "" {
			author = "Unknown"
		}

		// Send job to worker
		jobChan <- QuoteJob{
			Content: quoteContent,
			Author:  author,
			Row:     rowCount,
		}
	}

	// Close job channel and wait for workers to finish
	close(jobChan)
	wg.Wait()
	done <- true

	duration := time.Since(startTime)
	log.Printf("=== SUMMARY for %s ===", filePath)
	log.Printf("Processing time: %v", duration)
	log.Printf("Total rows read: %d", rowCount)
	log.Printf("Processed: %d quotes", p.stats.Processed.Load())
	log.Printf("Errors: %d", p.stats.Errors.Load())
	log.Printf("Skipped: %d", p.stats.Skipped.Load())
	log.Printf("Total attempted: %d", p.stats.Processed.Load()+p.stats.Errors.Load()+p.stats.Skipped.Load())

	if rowCount > 0 {
		log.Printf("Processing rate: %.2f quotes/second", float64(p.stats.Processed.Load())/duration.Seconds())
	}

	return nil
}

func main() {
	filePath := flag.String("file", "quotes.csv", "The path to the CSV file or ZIP file")
	configPath := flag.String("config", "configs/config.yml", "The path to the config file")
	dryRun := flag.Bool("dry-run", false, "If true, only parse and validate quotes without inserting to DB")
	workers := flag.Int("workers", 10, "Number of worker goroutines")
	batchSize := flag.Int("batch-size", 100, "Number of quotes to process in each batch")
	extractDir := flag.String("extract-dir", "", "Directory to extract ZIP files to (optional)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadWithPath(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dbManager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbManager.Close()

	// Create quote repository
	quoteRepo := repository.NewPostgreSQLQuoteRepository(dbManager.Postgres)

	// Initialize stats
	stats := &WorkerStats{
		Processed: &atomic.Uint32{},
		Errors:    &atomic.Uint32{},
		Skipped:   &atomic.Uint32{},
	}

	// Create processor
	processor := &CSVProcessor{
		quoteRepo: quoteRepo,
		dryRun:    *dryRun,
		workers:   *workers,
		batchSize: *batchSize,
		stats:     stats,
	}

	// Process file based on extension
	fileExt := strings.ToLower(filepath.Ext(*filePath))

	switch fileExt {
	case ".zip":
		log.Printf("Processing ZIP file: %s", *filePath)
		if err := processZipFile(*filePath, *extractDir, processor); err != nil {
			log.Fatalf("Failed to process ZIP file: %v", err)
		}
	case ".csv":
		log.Printf("Processing CSV file: %s", *filePath)
		if err := processor.ProcessFile(*filePath); err != nil {
			log.Fatalf("Failed to process CSV file: %v", err)
		}
	default:
		log.Fatalf("Unsupported file type: %s. Only .csv and .zip files are supported", fileExt)
	}

	// Final summary
	log.Printf("=== FINAL SUMMARY ===")
	log.Printf("Total processed: %d quotes", stats.Processed.Load())
	log.Printf("Total errors: %d", stats.Errors.Load())
	log.Printf("Total skipped: %d", stats.Skipped.Load())

	if *dryRun {
		log.Println("DRY RUN completed - no data was inserted to database")
	} else {
		log.Println("Quote import completed successfully!")
	}
}

// processZipFile extracts and processes CSV files from a ZIP archive
func processZipFile(zipPath, extractDir string, processor *CSVProcessor) error {
	// Open ZIP file
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer zipReader.Close()

	// Create temporary extraction directory if not specified
	if extractDir == "" {
		tempDir, err := os.MkdirTemp("", "quotes_extract_*")
		if err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)
		extractDir = tempDir
	} else {
		// Create extraction directory if it doesn't exist
		if err := os.MkdirAll(extractDir, 0755); err != nil {
			return fmt.Errorf("failed to create extraction directory: %w", err)
		}
	}

	log.Printf("Extracting ZIP to: %s", extractDir)

	// Extract and process CSV files
	csvFiles := []string{}
	for _, file := range zipReader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".csv") {
			// Extract CSV file
			extractedPath := filepath.Join(extractDir, filepath.Base(file.Name))
			if err := extractZipFile(file, extractedPath); err != nil {
				log.Printf("Warning: Failed to extract %s: %v", file.Name, err)
				continue
			}
			csvFiles = append(csvFiles, extractedPath)
			log.Printf("Extracted: %s -> %s", file.Name, extractedPath)
		}
	}

	if len(csvFiles) == 0 {
		return fmt.Errorf("no CSV files found in ZIP archive")
	}

	log.Printf("Found %d CSV files in ZIP", len(csvFiles))

	// Process each CSV file
	for i, csvFile := range csvFiles {
		log.Printf("Processing CSV file %d/%d: %s", i+1, len(csvFiles), filepath.Base(csvFile))
		if err := processor.ProcessFile(csvFile); err != nil {
			log.Printf("Warning: Failed to process %s: %v", csvFile, err)
			continue
		}
	}

	return nil
}

// extractZipFile extracts a single file from ZIP to destination
func extractZipFile(zipFile *zip.File, destPath string) error {
	// Open source file in ZIP
	srcFile, err := zipFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in ZIP: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy content
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// worker processes quotes from the job channel
func worker(ctx context.Context, id int, jobChan <-chan QuoteJob, quoteRepo repositories.QuoteRepository, dryRun bool, stats *WorkerStats, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Worker %d started", id)

	for job := range jobChan {
		// Create quote entity
		quote, err := entities.NewQuote(job.Content, job.Author)
		if err != nil {
			log.Printf("Worker %d: Failed to create quote entity at row %d: %v", id, job.Row, err)
			stats.Errors.Add(1)
			continue
		}

		if !dryRun {
			// Insert quote to database
			if err := quoteRepo.Create(ctx, quote); err != nil {
				log.Printf("Worker %d: Failed to insert quote at row %d: %v", id, job.Row, err)
				stats.Errors.Add(1)
				continue
			}
		}

		stats.Processed.Add(1)
	}

	log.Printf("Worker %d finished", id)
}

// progressReporter reports progress every few seconds
func progressReporter(stats *WorkerStats, done <-chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processed := stats.Processed.Load()
			errors := stats.Errors.Load()
			skipped := stats.Skipped.Load()
			total := processed + errors + skipped

			if total > 0 {
				log.Printf("Progress: %d processed, %d errors, %d skipped (total: %d)",
					processed, errors, skipped, total)
			}
		case <-done:
			return
		}
	}
}
