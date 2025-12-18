package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/book"
	"experiment/internal/handlers"
	"experiment/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *chi.Mux {
	bookStore := store.NewMemBookStore()
	h := handlers.NewHandler(bookStore)

	r := chi.NewRouter()
	r.Get("/books", h.ListBooks)
	r.Post("/books", h.CreateBook)
	r.Get("/books/{id}", h.GetBook)

	return r
}

func TestListBooks(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var books []book.Book
	err := json.NewDecoder(w.Body).Decode(&books)
	require.NoError(t, err)
	assert.Empty(t, books)
}

func TestCreateBook(t *testing.T) {
	r := setupTestRouter()

	// Test successful creation
	newBookReq := book.CreateBookRequest{
		Title:  "The Hitchhiker's Guide to the Galaxy",
		Author: "Douglas Adams",
		Pages:  193,
	}
	body, _ := json.Marshal(newBookReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var createdBook book.Book
	err := json.NewDecoder(w.Body).Decode(&createdBook)
	require.NoError(t, err)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, newBookReq.Title, createdBook.Title)
	assert.Equal(t, newBookReq.Author, createdBook.Author)
	assert.Equal(t, newBookReq.Pages, createdBook.Pages)

	// Test validation error
	invalidBookReq := book.CreateBookRequest{
		Title:  "",   // Too short
		Author: "DA", // Too short
		Pages:  0,    // Too small
	}
	body, _ = json.Marshal(invalidBookReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBook(t *testing.T) {
	bookStore := store.NewMemBookStore()
	handler := handlers.NewHandler(bookStore)

	r := chi.NewRouter()
	r.Get("/books/{id}", handler.GetBook)

	// Add a book to the store directly for testing
	existingBook := book.Book{
		ID:     "123e4567-e89b-12d3-a456-426614174000",
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}
	_ = bookStore.CreateBook(existingBook)

	// Test successful retrieval
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", existingBook.ID), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var retrievedBook book.Book
	err := json.NewDecoder(w.Body).Decode(&retrievedBook)
	require.NoError(t, err)
	assert.Equal(t, existingBook, retrievedBook)

	// Test book not found
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/books/non-existent-id", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Book not found")
}
