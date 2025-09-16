package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atdevten/peace/internal/infrastructure/config"
	"github.com/atdevten/peace/internal/infrastructure/database"
	"github.com/atdevten/peace/internal/interfaces/http/server"
)

func main() {
	// Parse command line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadWithPath(configPath)
	if err != nil {
		panic(err)
	}

	// Initialize database connections
	dbManager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		panic(err)
	}
	defer dbManager.Close()

	// Initialize HTTP server
	srv, err := server.NewHTTPServer(cfg)
	if err != nil {
		panic(err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Run(); err != nil {
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
