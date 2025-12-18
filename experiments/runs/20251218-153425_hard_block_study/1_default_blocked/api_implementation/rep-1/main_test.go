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
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	listBooks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var bookList []Book
	err := json.Unmarshal(rr.Body.Bytes(), &bookList)
	assert.NoError(t, err)
	assert.Empty(t, bookList)

	// Add a book and test again
	mu.Lock()
	books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()

	listBooks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	err = json.Unmarshal(rr.Body.Bytes(), &bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 1)
	assert.Equal(t, "Test Book", bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear books before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	// Valid request
	validReq := CreateBookRequest{Title: "New Book", Author: "Some Author", Pages: 150}
	reqBody, _ := json.Marshal(validReq)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	createBook(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var book Book
	err := json.Unmarshal(rr.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.NotEmpty(t, book.ID)
	assert.Equal(t, validReq.Title, book.Title)
	assert.Equal(t, validReq.Author, book.Author)
	assert.Equal(t, validReq.Pages, book.Pages)

	mu.RLock()
	assert.Len(t, books, 1)
	mu.RUnlock()

	// Invalid request body (non-json)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	createBook(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Validation error: Title too short
	invalidTitleReq := CreateBookRequest{Title: "ab", Author: "Author", Pages: 10}
	reqBody, _ = json.Marshal(invalidTitleReq)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	createBook(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Validation error: Pages less than 1
	invalidPagesReq := CreateBookRequest{Title: "Valid Title", Author: "Valid Author", Pages: 0}
	reqBody, _ = json.Marshal(invalidPagesReq)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	createBook(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Clear books before each test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	router := chi.NewRouter()
	router.Get("/books/{id}", getBook)

	// Book not found
	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Book not found")

	// Add a book and retrieve it
	existingBook := Book{ID: "abc-123", Title: "Existing Book", Author: "Existing Author", Pages: 200}
	mu.Lock()
	books["abc-123"] = existingBook
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books/abc-123", nil)
	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var book Book
	err := json.Unmarshal(rr.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, existingBook.ID, book.ID)
	assert.Equal(t, existingBook.Title, book.Title)
}
