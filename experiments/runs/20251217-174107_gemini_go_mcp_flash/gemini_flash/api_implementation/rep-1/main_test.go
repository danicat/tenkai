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
	// Initialize bookstore with some data
	mu.Lock()
	bookstore = make(map[string]Book)
	bookstore["1"] = Book{
		ID:    "1",
		Title:  "The Lord of the Rings",
		Author: "J.R.R. Tolkien",
		Pages:  1200,
	}
	bookstore["2"] = Book{
		ID:    "2",
		Title:  "The Hobbit",
		Author: "J.R.R. Tolkien",
		Pages:  300,
	}
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/books", listBooks)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var books []Book
	err := json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 2)
}

func TestCreateBook(t *testing.T) {
	mu.Lock()
	bookstore = make(map[string]Book)
	mu.Unlock()

	newBook := CreateBookRequest{
		Title:  "The Hitchhiker's Guide to the Galaxy",
		Author: "Douglas Adams",
		Pages:  193,
	}
	body, _ := json.Marshal(newBook)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Post("/books", createBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var book Book
	err := json.NewDecoder(rr.Body).Decode(&book)
	assert.NoError(t, err)
	assert.Equal(t, newBook.Title, book.Title)
	assert.Equal(t, newBook.Author, book.Author)
	assert.Equal(t, newBook.Pages, book.Pages)
	assert.NotEmpty(t, book.ID)

	mu.RLock()
	_, ok := bookstore[book.ID]
	mu.RUnlock()
	assert.True(t, ok)
}

func TestCreateBook_InvalidInput(t *testing.T) {
	mu.Lock()
	bookstore = make(map[string]Book)
	mu.Unlock()

	invalidBook := CreateBookRequest{
		Title:  "", // Invalid: minLength 3
		Author: "",
		Pages:  0, // Invalid: minimum 1
	}
	body, _ := json.Marshal(invalidBook)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Post("/books", createBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	// Initialize bookstore with some data
	mu.Lock()
	bookstore = make(map[string]Book)
	bookstore["123"] = Book{
		ID:    "123",
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/123", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var book Book
	err := json.NewDecoder(rr.Body).Decode(&book)
	assert.NoError(t, err)
	assert.Equal(t, "123", book.ID)
}

func TestGetBook_NotFound(t *testing.T) {
	mu.Lock()
	bookstore = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}