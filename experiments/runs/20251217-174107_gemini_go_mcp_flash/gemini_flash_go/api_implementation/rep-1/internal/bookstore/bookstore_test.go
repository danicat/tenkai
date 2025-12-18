package bookstore

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	s := NewServer()
	r := s.Routes()

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var books []Book
	err := json.NewDecoder(w.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Empty(t, books) // Initially, no books

	// Add a book and test again
	s.Store.CreateBook(context.Background(), CreateBookRequest{
		Title: "Test Book", Author: "Test Author", Pages: 100,
	})

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	books = []Book{}
	err = json.NewDecoder(w.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)
}

func TestCreateBook(t *testing.T) {
	s := NewServer()
	r := s.Routes()

	// Success case
	createReq := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var book Book
	err := json.NewDecoder(w.Body).Decode(&book)
	assert.NoError(t, err)
	assert.NotEmpty(t, book.ID)
	assert.Equal(t, "New Book", book.Title)

	// Validation error case (title too short)
	invalidReq := CreateBookRequest{Title: "AB", Author: "New Author", Pages: 200}
	reqBody, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Key: 'CreateBookRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag")

	// Invalid JSON body
	req = httptest.NewRequest("POST", "/books", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestGetBook(t *testing.T) {
	s := NewServer()
	r := s.Routes()

	// Create a book first
	createdBook, _ := s.Store.CreateBook(context.Background(), CreateBookRequest{
		Title: "Get This Book", Author: "Author A", Pages: 123,
	})

	// Success case
	req := httptest.NewRequest("GET", "/books/"+createdBook.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var fetchedBook Book
	err := json.NewDecoder(w.Body).Decode(&fetchedBook)
	assert.NoError(t, err)
	assert.Equal(t, createdBook.ID, fetchedBook.ID)
	assert.Equal(t, "Get This Book", fetchedBook.Title)

	// Not found case
	req = httptest.NewRequest("GET", "/books/nonexistent-id", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "book not found")

	// Empty ID (chi.URLParam will return empty string, handled by handler). This will hit /books for list, not get by ID.
	req = httptest.NewRequest("GET", "/books", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code) // Should hit list books and return empty array

	// Test with a router that ONLY has the GetBook route, to ensure 404 for missing ID.
	t.Run("empty id with specific router", func(t *testing.T) {
		isolatedRouter := chi.NewRouter()
		isolatedRouter.Get("/books/{id}", s.handleGetBook())

		req = httptest.NewRequest("GET", "/books/", nil)
		w = httptest.NewRecorder()
		isolatedRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code) // If only get by ID is defined, then a missing ID will be handled as 404 by chi
	})

}
