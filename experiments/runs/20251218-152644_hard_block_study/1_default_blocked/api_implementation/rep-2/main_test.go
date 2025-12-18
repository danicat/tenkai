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
	// Clear books map before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books", listBooks)

	// Add some dummy books
	book1 := Book{ID: "1", Title: "Book One", Author: "Author One"}
	book2 := Book{ID: "2", Title: "Book Two", Author: "Author Two"}
	mu.Lock()
	books[book1.ID] = book1
	books[book2.ID] = book2
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBooks []Book
	err := json.NewDecoder(rr.Body).Decode(&responseBooks)
	assert.NoError(t, err)

	assert.Len(t, responseBooks, 2)
	assert.Contains(t, responseBooks, book1)
	assert.Contains(t, responseBooks, book2)
}

func TestCreateBook(t *testing.T) {
	// Clear books map before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Post("/books", createBook)

	newBook := NewBook{Title: "New Book", Author: "New Author"}
	body, _ := json.Marshal(newBook)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdBook Book
	err := json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.NoError(t, err)

	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)

	// Verify book is in storage
	mu.RLock()
	defer mu.RUnlock()
	_, ok := books[createdBook.ID]
	assert.True(t, ok)
}

func TestCreateBook_InvalidInput(t *testing.T) {
	// Clear books map before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Missing title
	newBook := NewBook{Title: "", Author: "New Author"}
	body, _ := json.Marshal(newBook)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errorResponse map[string]string
	err := json.NewDecoder(rr.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse["error"], "Validation error")
}

func TestGetBook(t *testing.T) {
	// Clear books map before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Add a dummy book
	book1 := Book{ID: "123", Title: "Test Book", Author: "Test Author"}
	mu.Lock()
	books[book1.ID] = book1
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBook Book
	err := json.NewDecoder(rr.Body).Decode(&responseBook)
	assert.NoError(t, err)

	assert.Equal(t, book1, responseBook)
}

func TestGetBook_NotFound(t *testing.T) {
	// Clear books map before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var errorResponse map[string]string
	err := json.NewDecoder(rr.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Book not found", errorResponse["error"])
}
