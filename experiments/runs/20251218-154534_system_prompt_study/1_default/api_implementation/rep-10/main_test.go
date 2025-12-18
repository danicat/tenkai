package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	// Initialize router and handlers for testing
	r := chi.NewRouter()
	r.Get("/books", listBooks)

	// Clear existing books and add a test book
	storeMu.Lock()
	bookStore = make(map[string]Book)
	bookStore["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	storeMu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var books []Book
	err := json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Initialize router and handlers for testing
	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Clear existing books
	storeMu.Lock()
	bookStore = make(map[string]Book)
	storeMu.Unlock()

	// Valid book creation
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdBook Book
	err := json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.Equal(t, "New Book", createdBook.Title)
	assert.NotEmpty(t, createdBook.ID)

	// Check if the book is in the store
	storeMu.RLock()
	defer storeMu.RUnlock()
	_, ok := bookStore[createdBook.ID]
	assert.True(t, ok)

	// Invalid book creation (missing title)
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Initialize router and handlers for testing
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Clear existing books and add a test book
	storeMu.Lock()
	bookStore = make(map[string]Book)
	bookStore["test-123"] = Book{ID: "test-123", Title: "Fetch Book", Author: "Fetch Author", Pages: 300}
	storeMu.Unlock()

	// Get existing book
	req, _ := http.NewRequest("GET", "/books/test-123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var fetchedBook Book
	err := json.NewDecoder(rr.Body).Decode(&fetchedBook)
	assert.NoError(t, err)
	assert.Equal(t, "Fetch Book", fetchedBook.Title)

	// Get non-existent book
	req, _ = http.NewRequest("GET", "/books/non-existent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
