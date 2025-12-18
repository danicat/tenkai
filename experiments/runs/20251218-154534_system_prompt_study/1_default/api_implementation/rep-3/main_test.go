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
	// Clear the bookStore for a clean test environment
	storeMu.Lock()
	bookStore = make(map[string]Book)
	storeMu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	ListBooks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var books []Book
	err := json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Empty(t, books)

	// Add a book and test again
	book := Book{
		ID:     "test-id-123",
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}

	storeMu.Lock()
	bookStore[book.ID] = book
	storeMu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()

	ListBooks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	err = json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, book.Title, books[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear the bookStore for a clean test environment
	storeMu.Lock()
	bookStore = make(map[string]Book)
	storeMu.Unlock()

	// Valid book creation
	validBook := CreateBookRequest{
		Title:  "New Valid Book",
		Author: "John Doe",
		Pages:  250,
	}
	body, _ := json.Marshal(validBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateBook(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var createdBook Book
	err := json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, validBook.Title, createdBook.Title)
	assert.Equal(t, validBook.Author, createdBook.Author)
	assert.Equal(t, validBook.Pages, createdBook.Pages)

	// Invalid book creation - missing title
	invalidBook := CreateBookRequest{
		Author: "Jane Doe",
		Pages:  120,
	}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	CreateBook(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid book creation - pages less than 1
	invalidBook = CreateBookRequest{
		Title:  "Bad Book",
		Author: "Jane Doe",
		Pages:  0,
	}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	CreateBook(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Clear the bookStore for a clean test environment
	storeMu.Lock()
	bookStore = make(map[string]Book)
	storeMu.Unlock()

	// Add a book to retrieve
	existingBook := Book{
		ID:     "existing-id-456",
		Title:  "Existing Book",
		Author: "Existing Author",
		Pages:  300,
	}
	storeMu.Lock()
	bookStore[existingBook.ID] = existingBook
	storeMu.Unlock()

	// Get existing book
	req, _ := http.NewRequest("GET", "/books/"+existingBook.ID, nil)
	rr := httptest.NewRecorder()

	// Need a Chi router to properly handle URL parameters
	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var retrievedBook Book
	err := json.NewDecoder(rr.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook.ID, retrievedBook.ID)
	assert.Equal(t, existingBook.Title, retrievedBook.Title)

	// Get non-existent book
	req, _ = http.NewRequest("GET", "/books/non-existent-id", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
