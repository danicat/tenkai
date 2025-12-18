package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	// Initialize a new router for testing
	r := chi.NewRouter()
	r.Get("/books", ListBooks)

	// Clear existing books and add a test book
	mu.Lock()
	bookstore = make(map[string]Book)
	book1 := Book{ID: uuid.New().String(), Title: "Test Book 1", Author: "Author 1", Pages: 100}
	bookstore[book1.ID] = book1
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var books []Book
	err := json.NewDecoder(rr.Body).Decode(&books)
	assert.Nil(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, book1.Title, books[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Initialize a new router for testing
	r := chi.NewRouter()
	r.Post("/books", CreateBook)

	// Clear existing books
	mu.Lock()
	bookstore = make(map[string]Book)
	mu.Unlock()

	newBook := CreateBookRequest{
		Title:  "New Book Title",
		Author: "New Author",
		Pages:  150,
	}
	body, _ := json.Marshal(newBook)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdBook Book
	err := json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.Nil(t, err)
	assert.NotNil(t, createdBook.ID)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)
	assert.Equal(t, newBook.Pages, createdBook.Pages)

	mu.RLock()
	defer mu.RUnlock()
	assert.Contains(t, bookstore, createdBook.ID)
}

func TestCreateBook_InvalidInput(t *testing.T) {
	// Initialize a new router for testing
	r := chi.NewRouter()
	r.Post("/books", CreateBook)

	// Test case: Missing title
	invalidBook := CreateBookRequest{
		Author: "Invalid Author",
		Pages:  10,
	}
	body, _ := json.Marshal(invalidBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test case: Title too short
	invalidBook = CreateBookRequest{
		Title:  "ab",
		Author: "Invalid Author",
		Pages:  10,
	}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test case: Pages less than 1
	invalidBook = CreateBookRequest{
		Title:  "Valid Title",
		Author: "Valid Author",
		Pages:  0,
	}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Initialize a new router for testing
	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)

	// Clear existing books and add a test book
	mu.Lock()
	bookstore = make(map[string]Book)
	book1 := Book{ID: uuid.New().String(), Title: "Test Book 1", Author: "Author 1", Pages: 100}
	bookstore[book1.ID] = book1
	mu.Unlock()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", book1.ID), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedBook Book
	err := json.NewDecoder(rr.Body).Decode(&retrievedBook)
	assert.Nil(t, err)
	assert.Equal(t, book1.ID, retrievedBook.ID)
}

func TestGetBook_NotFound(t *testing.T) {
	// Initialize a new router for testing
	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)

	// Clear existing books
	mu.Lock()
	bookstore = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/nonexistent-id", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
