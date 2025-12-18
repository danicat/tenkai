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
	r := chi.NewRouter()
	r.Get("/books", listBooks)

	// Add some dummy books for testing
	mu.Lock()
	books = map[string]Book{
		"1": {ID: "1", Title: "Test Book 1", Author: "Author 1", Pages: 100},
		"2": {ID: "2", Title: "Test Book 2", Author: "Author 2", Pages: 200},
	}
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseBooks []Book
	json.NewDecoder(rr.Body).Decode(&responseBooks)
	assert.Len(t, responseBooks, 2)
	assert.Contains(t, responseBooks, Book{ID: "1", Title: "Test Book 1", Author: "Author 1", Pages: 100})
	assert.Contains(t, responseBooks, Book{ID: "2", Title: "Test Book 2", Author: "Author 2", Pages: 200})
}

func TestCreateBook(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Clear books for a clean test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	requestBody := CreateBookRequest{
		Title:  "New Book",
		Author: "New Author",
		Pages:  150,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var responseBook Book
	json.NewDecoder(rr.Body).Decode(&responseBook)
	assert.NotNil(t, responseBook.ID)
	assert.Equal(t, "New Book", responseBook.Title)
	assert.Equal(t, "New Author", responseBook.Author)
	assert.Equal(t, 150, responseBook.Pages)

	// Verify book is in store
	mu.RLock()
	defer mu.RUnlock()
	_, ok := books[responseBook.ID]
	assert.True(t, ok)
}

func TestCreateBookInvalidInput(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Missing title
	requestBody := `{"author": "Author", "pages": 100}`
	req, _ := http.NewRequest("POST", "/books", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid pages
	requestBody = `{"title": "Title", "author": "Author", "pages": 0}`
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Add a dummy book
	mu.Lock()
	books = map[string]Book{
		"123": {ID: "123", Title: "Specific Book", Author: "Specific Author", Pages: 300},
	}
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseBook Book
	json.NewDecoder(rr.Body).Decode(&responseBook)
	assert.Equal(t, Book{ID: "123", Title: "Specific Book", Author: "Specific Author", Pages: 300}, responseBook)
}

func TestGetBookNotFound(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Clear books for a clean test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
