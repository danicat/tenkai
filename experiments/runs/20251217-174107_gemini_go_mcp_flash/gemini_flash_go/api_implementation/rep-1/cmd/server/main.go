package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"experiment/internal/bookstore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Printf("server terminated: %v", err)
		os.Exit(1)
	}
	log.Println("server gracefully stopped")
}

func run(ctx context.Context) error {
	srv := bookstore.NewServer()
	router := srv.Routes()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v", httpServer.Addr, err)
		}
	}()

	// Wait for context cancellation (e.g., interrupt signal)
	<-ctx.Done()

	// Shut down server gracefully
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	log.Println("server shutting down...")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
