package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	// Initialize the router and clear books for a clean test state
	r := chi.NewRouter()
	r.Get("/books", listBooks)
	books = make(map[string]Book)

	// Add a dummy book
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response []Book
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Len(t, response, 1)
	assert.Equal(t, "Test Book", response[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Initialize the router and clear books for a clean test state
	r := chi.NewRouter()
	r.Post("/books", createBook)
	books = make(map[string]Book)
	validate = validator.New() // Re-initialize validator

	// Valid request
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response Book
	json.NewDecoder(rr.Body).Decode(&response)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "New Book", response.Title)
	assert.Equal(t, 1, len(books)) // Check if book was added to storage

	// Invalid request (missing title)
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, 1, len(books)) // Should not have added the invalid book

	// Invalid request (pages too low)
	invalidBook2 := CreateBookRequest{Title: "Book 3", Author: "Author 3", Pages: 0}
	body, _ = json.Marshal(invalidBook2)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, 1, len(books)) // Should not have added the invalid book

}

func TestGetBook(t *testing.T) {
	// Initialize the router and clear books for a clean test state
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)
	books = make(map[string]Book)

	// Add a dummy book
	existingBook := Book{ID: "abc-123", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	books["abc-123"] = existingBook

	// Valid request
	req, _ := http.NewRequest("GET", "/books/abc-123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response Book
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "Existing Book", response.Title)

	// Not found request
	req, _ = http.NewRequest("GET", "/books/non-existent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Book not found")
}
