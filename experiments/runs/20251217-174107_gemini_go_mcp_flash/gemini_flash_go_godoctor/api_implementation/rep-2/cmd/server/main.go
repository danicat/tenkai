package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"experiment/internal/handlers"
	"experiment/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	bookStore := store.NewMemBookStore()
	handler := handlers.NewHandler(bookStore)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/books", handler.ListBooks)
	r.Post("/books", handler.CreateBook)
	r.Get("/books/{id}", handler.GetBook)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	}()

	log.Printf("server listening on %s", server.Addr)
	return server.ListenAndServe()
}
