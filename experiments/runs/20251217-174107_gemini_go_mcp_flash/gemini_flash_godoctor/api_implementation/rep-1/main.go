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

// Book represents a book in the store
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// CreateBookRequest represents the request body for creating a book
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// In-memory store for books
var (
	bookStore = make(map[string]Book)
	mu        sync.RWMutex
	validate  *validator.Validate
)

func init() {
	validate = validator.New()
}

// listBooksHandler returns a list of all books
func listBooksHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	books := make([]Book, 0, len(bookStore))
	for _, book := range bookStore {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// createBookHandler creates a new book
func createBookHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	mu.Lock()
	defer mu.Unlock()
	bookStore[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// getBookHandler returns a book by its ID
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.RLock()
	defer mu.RUnlock()
	book, ok := bookStore[id]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", listBooksHandler)
	r.Post("/books", createBookHandler)
	r.Get("/books/{id}", getBookHandler)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
