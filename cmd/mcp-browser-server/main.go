package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/periplon/bract/internal/browser"
	"github.com/periplon/bract/internal/config"
	"github.com/periplon/bract/internal/handler"
	"github.com/periplon/bract/internal/mcp"
	"github.com/periplon/bract/internal/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create browser client for Chrome extension communication
	browserClient := browser.NewClient(cfg.WebSocket)

	// Start WebSocket server for Chrome extension
	wsServer := websocket.NewServer(cfg.WebSocket.Port, browserClient)
	go func() {
		if err := wsServer.Start(ctx); err != nil {
			log.Printf("WebSocket server error: %v", err)
		}
	}()

	// Create tool handler with browser client
	toolHandler := handler.NewBrowserHandler(browserClient)

	// Create and configure MCP server
	mcpServer := mcp.NewServer(cfg.Server.Name, cfg.Server.Version, toolHandler)

	// Start MCP server in separate goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- mcpServer.Start()
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		log.Println("Shutting down...")
		cancel()
	case err := <-errChan:
		if err != nil {
			log.Fatalf("MCP server error: %v", err)
		}
	}

	fmt.Println("Server stopped")
}