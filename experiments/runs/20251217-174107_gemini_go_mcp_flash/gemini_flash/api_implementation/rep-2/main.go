package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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

// In-memory storage for books
var (
	books = make(map[string]Book)
	mu    sync.RWMutex // Mutex to protect access to the books map
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func main() {
	r := chi.NewRouter()

	r.Get("/books", ListBooks)
	r.Post("/books", CreateBook)
	r.Get("/books/{id}", GetBook)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}

// ListBooks handles GET /books requests.
func ListBooks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookList)
}

// CreateBook handles POST /books requests.
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

	mu.Lock()
	defer mu.Unlock()

	id := uuid.New().String()
	book := Book{
		ID:     id,
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}
	books[id] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// GetBook handles GET /books/{id} requests.
func GetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.RLock()
	defer mu.RUnlock()

	book, found := books[id]
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
