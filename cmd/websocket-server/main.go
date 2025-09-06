package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atdevten/peace/internal/infrastructure/config"
	"github.com/atdevten/peace/internal/interfaces/websocket/server"
)

func main() {
	// Parse command line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadWithPath(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize WebSocket server
	wsServer, err := server.NewWebSocketServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create WebSocket server: %v", err)
	}

	// Start WebSocket server
	if err := wsServer.Run(); err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}

	log.Printf("WebSocket server started")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down WebSocket server...")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := wsServer.Shutdown(ctx); err != nil {
		log.Fatalf("WebSocket server forced to shutdown: %v", err)
	}

	log.Println("WebSocket server exited")
}
