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
	// Clear books before each test
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Get("/books", ListBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var bookList []Book
	err := json.Unmarshal(rec.Body.Bytes(), &bookList)
	assert.NoError(t, err)
	assert.Empty(t, bookList)

	// Add a book and test again
	mu.Lock()
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &bookList)
	assert.Len(t, bookList, 1)
	assert.Equal(t, "Test Book", bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear books before each test
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Post("/books", CreateBook)

	// Valid request
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var createdBook Book
	err := json.Unmarshal(rec.Body.Bytes(), &createdBook)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, "New Book", createdBook.Title)

	mu.RLock()
	assert.Contains(t, books, createdBook.ID)
	mu.RUnlock()

	// Invalid request - missing title
	invalidBook := CreateBookRequest{Title: "", Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Invalid request - pages too low
	invalidBook = CreateBookRequest{Title: "Valid Title", Author: "Valid Author", Pages: 0}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Invalid request - non-json body
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString("not a json"))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetBook(t *testing.T) {
	// Clear books before each test
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)

	// Add a book for testing
	mu.Lock()
	books["test-id"] = Book{ID: "test-id", Title: "Fetched Book", Author: "Fetched Author", Pages: 300}
	mu.Unlock()

	// Get existing book
	req, _ := http.NewRequest("GET", "/books/test-id", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var fetchedBook Book
	err := json.Unmarshal(rec.Body.Bytes(), &fetchedBook)
	assert.NoError(t, err)
	assert.Equal(t, "test-id", fetchedBook.ID)
	assert.Equal(t, "Fetched Book", fetchedBook.Title)

	// Get non-existing book
	req, _ = http.NewRequest("GET", "/books/non-existent-id", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
