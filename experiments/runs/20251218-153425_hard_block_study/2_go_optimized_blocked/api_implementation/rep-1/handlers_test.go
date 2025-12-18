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
	bookStore := NewBookStore()
	h := &Handler{BookStore: bookStore}

	// Add some dummy books
	bookStore.AddBook(Book{ID: "1", Title: "Book 1", Author: "Author 1", Pages: 100})
	bookStore.AddBook(Book{ID: "2", Title: "Book 2", Author: "Author 2", Pages: 200})

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	h.ListBooks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var books []Book
	err := json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 2)
}

func TestCreateBook(t *testing.T) {
	bookStore := NewBookStore()
	h := &Handler{BookStore: bookStore}

	// Valid book creation
	validBook := CreateBookRequest{
		Title:  "New Book",
		Author: "New Author",
		Pages:  150,
	}
	body, _ := json.Marshal(validBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateBook(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var createdBook Book
	err := json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.Equal(t, validBook.Title, createdBook.Title)
	assert.NotEmpty(t, createdBook.ID)

	// Invalid book creation (missing title)
	invalidBook := CreateBookRequest{
		Author: "Invalid Author",
		Pages:  10,
	}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	h.CreateBook(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid book creation (title too short)
	shortTitleBook := CreateBookRequest{
		Title:  "ab",
		Author: "Author",
		Pages:  10,
	}
	body, _ = json.Marshal(shortTitleBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	h.CreateBook(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	bookStore := NewBookStore()
	h := &Handler{BookStore: bookStore}

	r := chi.NewRouter()
	r.Get("/books/{id}", h.GetBook)

	// Add a dummy book
	existingBook := Book{ID: "123", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	bookStore.AddBook(existingBook)

	// Get existing book
	req, _ := http.NewRequest("GET", "/books/123", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var retrievedBook Book
	err := json.NewDecoder(rr.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook, retrievedBook)

	// Get non-existent book
	req, _ = http.NewRequest("GET", "/books/nonexistent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
