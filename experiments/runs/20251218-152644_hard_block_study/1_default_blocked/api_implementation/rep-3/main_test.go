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
	// Initialize validator and clear books map for clean test
	validate = validator.New()
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Get("/books", listBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	var bookList []Book
	json.NewDecoder(rw.Body).Decode(&bookList)
	assert.Empty(t, bookList)

	// Add a book and test again
	mu.Lock()
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	json.NewDecoder(rw.Body).Decode(&bookList)
	assert.Len(t, bookList, 1)
	assert.Equal(t, "Test Book", bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Initialize validator and clear books map for clean test
	validate = validator.New()
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Valid request
	validBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(validBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusCreated, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	var createdBook Book
	json.NewDecoder(rw.Body).Decode(&createdBook)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, validBook.Title, createdBook.Title)
	assert.Equal(t, validBook.Author, createdBook.Author)
	assert.Equal(t, validBook.Pages, createdBook.Pages)

	mu.RLock()
	_, ok := books[createdBook.ID]
	mu.RUnlock()
	assert.True(t, ok)

	// Invalid request - missing title
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)

	// Invalid request - pages too low
	invalidBook = CreateBookRequest{Title: "Bad Pages", Author: "Author", Pages: 0}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

func TestGetBook(t *testing.T) {
	// Initialize validator and clear books map for clean test
	validate = validator.New()
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Book not found
	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Code)

	// Add a book and retrieve it
	mu.Lock()
	existingBook := Book{ID: "123", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	books["123"] = existingBook
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books/123", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	var retrievedBook Book
	json.NewDecoder(rw.Body).Decode(&retrievedBook)
	assert.Equal(t, existingBook, retrievedBook)
}
