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
	// Clear existing books for a clean test environment
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books", ListBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var bookList []Book
	err := json.NewDecoder(rr.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Empty(t, bookList)

	// Add a book and test again
	mu.Lock()
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	mu.Unlock()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	bookList = []Book{}
	err = json.NewDecoder(rr.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 1)
	assert.Equal(t, "Test Book", bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear existing books for a clean test environment
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Post("/books", CreateBook)

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
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.NotEmpty(t, createdBook.ID)

	// Invalid book creation (missing title)
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)

	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid book creation (pages < 1)
	invalidBook2 := CreateBookRequest{Title: "Title", Author: "Author", Pages: 0}
	body, _ = json.Marshal(invalidBook2)

	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Clear existing books for a clean test environment
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)

	// Book not found
	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Add a book and retrieve it
	existingBook := Book{ID: "2", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	mu.Lock()
	books[existingBook.ID] = existingBook
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books/2", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedBook Book
	err := json.NewDecoder(rr.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook.Title, retrievedBook.Title)
}
