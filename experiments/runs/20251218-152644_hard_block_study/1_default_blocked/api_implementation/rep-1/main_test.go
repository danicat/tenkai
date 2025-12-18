package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	// Initialize an empty book store for testing
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books", listBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var bookList []Book
	err := json.NewDecoder(w.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Empty(t, bookList)

	// Add a book and test again
	mu.Lock()
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	mu.Unlock()

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err = json.NewDecoder(w.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 1)
	assert.Equal(t, "Test Book", bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Initialize an empty book store for testing
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Valid book creation
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdBook Book
	err := json.NewDecoder(w.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.NotEmpty(t, createdBook.ID)

	// Invalid book creation - missing title
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse map[string]string
	err = json.NewDecoder(w.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse["error"], "Title") // Check for validation error on Title

	// Invalid book creation - pages less than 1
	invalidBook = CreateBookRequest{Title: "Book", Author: "Author", Pages: 0}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	err = json.NewDecoder(w.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse["error"], "Pages") // Check for validation error on Pages
}

func TestGetBook(t *testing.T) {
	// Initialize an empty book store for testing
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Book not found
	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Add a book and retrieve it
	mu.Lock()
	existingBook := Book{ID: "abc-123", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	books["abc-123"] = existingBook
	mu.Unlock()

	req, _ = http.NewRequest("GET", fmt.Sprintf("/books/%s", existingBook.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedBook Book
	err := json.NewDecoder(w.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook.ID, retrievedBook.ID)
	assert.Equal(t, existingBook.Title, retrievedBook.Title)
}
