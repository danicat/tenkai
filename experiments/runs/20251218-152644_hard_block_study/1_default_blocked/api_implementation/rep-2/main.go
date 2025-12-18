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

// Book represents a book in the bookstore
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

// NewBook represents a new book to be created
type NewBook struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

var (
	books    = make(map[string]Book)
	mu       sync.RWMutex
	validate = validator.New()
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", listBooks)
	r.Post("/books", createBook)
	r.Get("/books/{id}", getBook)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}

	writeJSONResponse(w, http.StatusOK, bookList)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook NewBook
	if err := readJSONRequest(r, &newBook); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validate.Struct(newBook); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", err.Error()))
		return
	}

	book := Book{
		ID:     uuid.New().String(),
		Title:  newBook.Title,
		Author: newBook.Author,
	}

	mu.Lock()
	defer mu.Unlock()
	books[book.ID] = book

	writeJSONResponse(w, http.StatusCreated, book)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	mu.RLock()
	defer mu.RUnlock()

	book, ok := books[id]
	if !ok {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, book)
}

// Helper functions for JSON responses and errors
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	writeJSONResponse(w, status, map[string]string{"error": message})
}

func readJSONRequest(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}
