package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Book represents a book in the bookstore.
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// CreateBookRequest represents the request body for creating a new book.
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// Server holds the application's dependencies.
type Server struct {
	books     map[string]Book
	validator *validator.Validate
	mu        sync.RWMutex
}

// NewServer creates and initializes a new Server.
func NewServer() *Server {
	return &Server{
		books:     make(map[string]Book),
		validator: validator.New(),
	}
}

// handleListBooks handles GET /books requests.
func (s *Server) handleListBooks(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// handleCreateBook handles POST /books requests.
func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	s.mu.Lock()
	s.books[book.ID] = book
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// handleGetBook handles GET /books/{id} requests.
func (s *Server) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	s.mu.RLock()
	book, ok := s.books[id]
	s.mu.RUnlock()

	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// setupRoutes sets up the HTTP routes.
func (s *Server) setupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/books", s.handleListBooks)
	r.Post("/books", s.handleCreateBook)
	r.Get("/books/{id}", s.handleGetBook)

	return r
}

func run(ctx context.Context, stdout *os.File) error {
	srv := NewServer()
	router := srv.setupRoutes()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", httpServer.Addr, err)
		}
	}()

	// Wait for context cancellation (e.g., interrupt signal)
	<-ctx.Done()

	// Shut down server gracefully
	log.Println("Server shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server gracefully stopped.")
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, os.Stdout); err != nil {

		fmt.Fprintf(os.Stderr, "Server error: %s\n", err)

		os.Exit(1)

	}

}
