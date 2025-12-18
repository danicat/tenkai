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
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

// CreateBookRequest represents the request body for creating a new book
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

var (
	books    = make(map[string]Book)
	booksMux sync.Mutex
	validate = validator.New()
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", ListBooks)
	r.Post("/books", CreateBook)
	r.Get("/books/{id}", GetBook)

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

// ListBooks handles GET /books requests
func ListBooks(w http.ResponseWriter, r *http.Request) {
	booksMux.Lock()
	defer booksMux.Unlock()

	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookList)
}

// CreateBook handles POST /books requests
func CreateBook(w http.ResponseWriter, r *http.Request) {
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

	booksMux.Lock()
	defer booksMux.Unlock()
	books[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// GetBook handles GET /books/{id} requests
func GetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	booksMux.Lock()
	defer booksMux.Unlock()

	book, found := books[id]
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
