package main

import (
	"context"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
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

func main() {
	filePath := flag.String("file", "quotes.csv", "The path to the CSV file")
	configPath := flag.String("config", "configs/config.yml", "The path to the config file")
	dryRun := flag.Bool("dry-run", false, "If true, only parse and validate quotes without inserting to DB")
	workers := flag.Int("workers", 10, "Number of worker goroutines")
	batchSize := flag.Int("batch-size", 100, "Number of quotes to process in each batch")
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

	// Open CSV file
	csvFile, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1

	ctx := context.Background()
	stats := &WorkerStats{
		Processed: &atomic.Uint32{},
		Errors:    &atomic.Uint32{},
		Skipped:   &atomic.Uint32{},
	}

	log.Printf("Starting to process quotes from %s", *filePath)
	log.Printf("Using %d workers with batch size %d", *workers, *batchSize)
	if *dryRun {
		log.Println("DRY RUN mode enabled - no data will be inserted")
	}

	// Create job channel and start workers
	jobChan := make(chan QuoteJob, *batchSize)
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(ctx, i+1, jobChan, quoteRepo, *dryRun, stats, &wg)
	}

	// Start progress reporter
	done := make(chan bool)
	go progressReporter(stats, done)

	// Skip header row
	if _, err := csvReader.Read(); err != nil {
		log.Fatalf("Failed to read header: %v", err)
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
			stats.Errors.Add(1)
			continue
		}

		rowCount++

		// Parse CSV record (quote, author, category)
		if len(record) < 2 {
			log.Printf("Skipping invalid record at row %d (insufficient columns)", rowCount)
			stats.Skipped.Add(1)
			continue
		}

		quoteContent := strings.TrimSpace(record[0])
		author := strings.TrimSpace(record[1])

		// Skip empty quotes
		if quoteContent == "" {
			log.Printf("Skipping empty quote at row %d", rowCount)
			stats.Skipped.Add(1)
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
	log.Printf("=== SUMMARY ===")
	log.Printf("Processing time: %v", duration)
	log.Printf("Total rows read: %d", rowCount)
	log.Printf("Processed: %d quotes", stats.Processed.Load())
	log.Printf("Errors: %d", stats.Errors.Load())
	log.Printf("Skipped: %d", stats.Skipped.Load())
	log.Printf("Total attempted: %d", stats.Processed.Load()+stats.Errors.Load()+stats.Skipped.Load())

	if rowCount > 0 {
		log.Printf("Processing rate: %.2f quotes/second", float64(stats.Processed.Load())/duration.Seconds())
	}

	if *dryRun {
		log.Println("DRY RUN completed - no data was inserted to database")
	} else {
		log.Println("Quote import completed successfully!")
	}
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
