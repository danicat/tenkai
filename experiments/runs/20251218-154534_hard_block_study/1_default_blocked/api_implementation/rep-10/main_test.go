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
	router := chi.NewRouter()
	router.Get("/books", listBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var books []Book
	json.NewDecoder(rr.Body).Decode(&books)
	assert.Empty(t, books)

	// Add a book and test again
	bookStore["1"] = Book{ID: "1", Title: "Test Book 1", Author: "Test Author 1", Pages: 100}
	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.NewDecoder(rr.Body).Decode(&books)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book 1", books[0].Title)
	// Clean up
	delete(bookStore, "1")
}

func TestCreateBook(t *testing.T) {
	router := chi.NewRouter()
	router.Post("/books", createBook)

	// Valid book creation
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var book Book
	json.NewDecoder(rr.Body).Decode(&book)
	assert.Equal(t, "New Book", book.Title)
	assert.NotEmpty(t, book.ID)

	// Check if book is in store
	storeMu.RLock()
	_, found := bookStore[book.ID]
	storeMu.RUnlock()
	assert.True(t, found)
	// Clean up
	delete(bookStore, book.ID)

	// Invalid book creation - missing title
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid book creation - pages < 1
	invalidBook = CreateBookRequest{Title: "Bad Pages", Author: "Author", Pages: 0}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBookByID(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/books/{id}", getBookByID)

	// Add a book for testing
	existingBook := Book{ID: "test-id-123", Title: "Existing Book", Author: "Existing Author", Pages: 150}
	storeMu.Lock()
	bookStore[existingBook.ID] = existingBook
	storeMu.Unlock()

	// Get existing book
	req, _ := http.NewRequest("GET", "/books/test-id-123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var book Book
	json.NewDecoder(rr.Body).Decode(&book)
	assert.Equal(t, existingBook.Title, book.Title)

	// Get non-existent book
	req, _ = http.NewRequest("GET", "/books/non-existent-id", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Clean up
	storeMu.Lock()
	delete(bookStore, existingBook.ID)
	storeMu.Unlock()
}
