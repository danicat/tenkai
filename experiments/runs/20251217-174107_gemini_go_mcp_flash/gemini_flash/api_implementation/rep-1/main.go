package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Book represents a book in the bookstore
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// CreateBookRequest represents the request body for creating a new book
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

var (
	bookstore = make(map[string]Book)
	mu        sync.RWMutex
	validate  = validator.New()
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", listBooks)
	r.Post("/books", createBook)
	r.Get("/books/{id}", getBook)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	books := make([]Book, 0, len(bookstore))
	for _, book := range bookstore {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
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
	bookstore[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.RLock()
	defer mu.RUnlock()
	book, ok := bookstore[id]
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
