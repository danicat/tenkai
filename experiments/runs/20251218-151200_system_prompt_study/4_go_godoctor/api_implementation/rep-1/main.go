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

// Book represents a book in the bookstore.
type Book struct {
	ID    string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages int    `json:"pages" validate:"min=1"`
}

// CreateBookRequest represents the request body for creating a new book.
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages int    `json:"pages" validate:"min=1"`
}

// In-memory storage for books
var (
	books = make(map[string]Book)
	mu    sync.RWMutex
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", listBooks)
	r.Post("/books", createBook)
	r.Get("/books/{id}", getBook)

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

// listBooks handles GET /books requests.
func listBooks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookList)
}

// createBook handles POST /books requests.
func createBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range ve {
				errors[e.Field()] = fmt.Sprintf("%s validation failed on the '%s' tag", e.Field(), e.Tag())
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookID := uuid.New().String()
	newBook := Book{
		ID:    bookID,
		Title:  req.Title,
		Author: req.Author,
		Pages: req.Pages,
	}

	mu.Lock()
	defer mu.Unlock()
	books[bookID] = newBook

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// getBook handles GET /books/{id} requests.
func getBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")

	mu.RLock()
	defer mu.RUnlock()

	book, found := books[bookID]
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
