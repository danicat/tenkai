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
	// Clear books for a clean slate
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books", listBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBooks []Book
	err := json.Unmarshal(rr.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)
	assert.Empty(t, responseBooks)

	// Add a book and test again
	mu.Lock()
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)
	assert.Len(t, responseBooks, 1)
	assert.Equal(t, "Test Book", responseBooks[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear books for a clean slate
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Test with valid input
	validBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(validBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdBook Book
	err := json.Unmarshal(rr.Body.Bytes(), &createdBook)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, validBook.Title, createdBook.Title)
	assert.Equal(t, validBook.Author, createdBook.Author)
	assert.Equal(t, validBook.Pages, createdBook.Pages)

	mu.RLock()
	assert.Contains(t, books, createdBook.ID)
	mu.RUnlock()

	// Test with invalid input (missing title)
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test with invalid input (pages < 1)
	invalidBook = CreateBookRequest{Title: "Invalid Book", Author: "Invalid Author", Pages: 0}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Clear books for a clean slate
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Add a book for testing
	existingBook := Book{ID: "test-123", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	mu.Lock()
	books[existingBook.ID] = existingBook
	mu.Unlock()

	// Test getting an existing book
	req, _ := http.NewRequest("GET", "/books/test-123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var fetchedBook Book
	err := json.Unmarshal(rr.Body.Bytes(), &fetchedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook.ID, fetchedBook.ID)
	assert.Equal(t, existingBook.Title, fetchedBook.Title)

	// Test getting a non-existent book
	req, _ = http.NewRequest("GET", "/books/non-existent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Book not found")
}
