package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
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

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Message string `json:"message"`
}

var (
	bookStore = make(map[string]Book)
	storeMu   sync.RWMutex
	validate  *validator.Validate
)

func init() {
	validate = validator.New()
}

func main() {
	router := chi.NewRouter()

	router.Get("/books", listBooks)
	router.Post("/books", createBook)
	router.Get("/books/{id}", getBookByID)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

// listBooks handles GET /books
func listBooks(w http.ResponseWriter, r *http.Request) {
	storeMu.RLock()
	defer storeMu.RUnlock()

	books := make([]Book, 0, len(bookStore))
	for _, book := range bookStore {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// createBook handles POST /books
func createBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	storeMu.Lock()
	defer storeMu.Unlock()
	bookStore[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// getBookByID handles GET /books/{id}
func getBookByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	storeMu.RLock()
	defer storeMu.RUnlock()

	book, found := bookStore[id]
	if !found {
		sendJSONError(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
