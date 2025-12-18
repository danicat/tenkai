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
	// Initialize the router and clear any existing books
	r := chi.NewRouter()
	r.Get("/books", listBooks)
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	// Create a sample book
	sampleBook := Book{
		ID:     "123",
		Title:  "The Go Programming Language",
		Author: "Alan A. A. Donovan",
		Pages:  400,
	}
	mu.Lock()
	books[sampleBook.ID] = sampleBook
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var responseBooks []Book
	err := json.NewDecoder(rec.Body).Decode(&responseBooks)
	assert.NoError(t, err)
	assert.Len(t, responseBooks, 1)
	assert.Equal(t, sampleBook.Title, responseBooks[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Initialize the router and clear any existing books
	r := chi.NewRouter()
	r.Post("/books", createBook)
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	requestBody := CreateBookRequest{
		Title:  "New Book Title",
		Author: "New Author",
		Pages:  250,
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var responseBook Book
	err := json.NewDecoder(rec.Body).Decode(&responseBook)
	assert.NoError(t, err)
	assert.NotNil(t, responseBook.ID)
	assert.Equal(t, requestBody.Title, responseBook.Title)
	assert.Equal(t, requestBody.Author, responseBook.Author)
	assert.Equal(t, requestBody.Pages, responseBook.Pages)

	// Verify the book is actually stored
	mu.RLock()
	defer mu.RUnlock()
	_, ok := books[responseBook.ID]
	assert.True(t, ok)
}

func TestCreateBook_InvalidInput(t *testing.T) {
	// Initialize the router and clear any existing books
	r := chi.NewRouter()
	r.Post("/books", createBook)
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	// Test case: Missing title
	invalidRequestBody := `{"author": "Test Author", "pages": 100}`
	req, _ := http.NewRequest("POST", "/books", bytes.NewBufferString(invalidRequestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test case: Title too short
	invalidRequestBody = `{"title": "ab", "author": "Test Author", "pages": 100}`
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString(invalidRequestBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test case: Pages less than 1
	invalidRequestBody = `{"title": "Valid Title", "author": "Test Author", "pages": 0}`
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString(invalidRequestBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetBook(t *testing.T) {
	// Initialize the router and clear any existing books
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	// Create a sample book
	sampleBook := Book{
		ID:     "test-id",
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  150,
	}
	mu.Lock()
	books[sampleBook.ID] = sampleBook
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/test-id", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var responseBook Book
	err := json.NewDecoder(rec.Body).Decode(&responseBook)
	assert.NoError(t, err)
	assert.Equal(t, sampleBook.ID, responseBook.ID)
	assert.Equal(t, sampleBook.Title, responseBook.Title)
}

func TestGetBook_NotFound(t *testing.T) {
	// Initialize the router and clear any existing books
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/nonexistent-id", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
