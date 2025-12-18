package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

// App represents the application state.
type App struct {
	books     map[string]Book
	validator *validator.Validate
	mu        sync.RWMutex
}

func NewApp() *App {
	return &App{
		books:     make(map[string]Book),
		validator: validator.New(),
	}
}

func (a *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", a.listBooks)
	r.Post("/books", a.createBook)
	r.Get("/books/{id}", a.getBook)

	return r
}

func main() {
	app := NewApp()

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (a *App) listBooks(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	books := make([]Book, 0, len(a.books))
	for _, book := range a.books {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *App) createBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := a.validator.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	book := Book{
		ID:     id,
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.books[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *App) getBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	a.mu.RLock()
	defer a.mu.RUnlock()

	book, ok := a.books[id]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
