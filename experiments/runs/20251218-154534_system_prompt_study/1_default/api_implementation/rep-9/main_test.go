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
	// Clear existing books and add some for testing
	books = make(map[string]Book)
	books["1"] = Book{ID: "1", Title: "Test Book 1", Author: "Author 1", Pages: 100}
	books["2"] = Book{ID: "2", Title: "Test Book 2", Author: "Author 2", Pages: 200}

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	ListBooks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var bookList []Book
	err := json.NewDecoder(rr.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 2)
	assert.Contains(t, bookList, books["1"])
	assert.Contains(t, bookList, books["2"])
}

func TestCreateBook(t *testing.T) {
	// Clear existing books
	books = make(map[string]Book)

	// Valid request
	validReq := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 150}
	body, _ := json.Marshal(validReq)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateBook(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdBook Book
	err := json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, validReq.Title, createdBook.Title)
	assert.Equal(t, validReq.Author, createdBook.Author)
	assert.Equal(t, validReq.Pages, createdBook.Pages)

	assert.Len(t, books, 1)
	_, found := books[createdBook.ID]
	assert.True(t, found)

	// Invalid request - missing title
	invalidReq := CreateBookRequest{Author: "New Author", Pages: 150}
	body, _ = json.Marshal(invalidReq)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	CreateBook(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid request - title too short
	shortTitleReq := CreateBookRequest{Title: "ab", Author: "New Author", Pages: 150}
	body, _ = json.Marshal(shortTitleReq)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	CreateBook(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Clear existing books and add one for testing
	books = make(map[string]Book)
	existingBook := Book{ID: "3", Title: "Existing Book", Author: "Author 3", Pages: 300}
	books["3"] = existingBook

	// Existing book
	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)

	req, _ := http.NewRequest("GET", "/books/3", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedBook Book
	err := json.NewDecoder(rr.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook, retrievedBook)

	// Non-existent book
	req, _ = http.NewRequest("GET", "/books/nonexistent", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
