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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

func main() {
	if err := run(); err != nil {
		log.Printf("server failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Configuration

	// For simplicity, we'll hardcode the port for this example.
	const port = ":8080"

	// =========================================================================
	// Dependencies

	store := bookstore.NewInMemBookStore()
	validate := validator.New()

	h := &bookstore.Handlers{
		Store:    store,
		Validate: validate,
	}

	// =========================================================================
	// Router Setup

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", h.ListBooks)
	r.Post("/books", h.CreateBook)
	r.Get("/books/{id}", h.GetBook)

	// =========================================================================
	// Start Server

	httpServer := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main: API listening on %s", port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		log.Println("main: Start shutdown")

		// Give outstanding requests a minute to complete.
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
