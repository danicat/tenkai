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
	router := chi.NewRouter()
	bookStore := NewBookStore()
	router.Get("/books", ListBooksHandler(bookStore))

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var books []Book
	err := json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Empty(t, books)

	// Add a book and test again
	book1 := Book{ID: uuid.New().String(), Title: "Test Book 1", Author: "Author 1", Pages: 100}
	bookStore.CreateBook(book1)

	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, book1.Title, books[0].Title)
}

func TestCreateBook(t *testing.T) {
	router := chi.NewRouter()
	bookStore := NewBookStore()
	router.Post("/books", CreateBookHandler(bookStore))

	createReq := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var book Book
	err := json.NewDecoder(rr.Body).Decode(&book)
	assert.NoError(t, err)
	assert.NotEmpty(t, book.ID)
	assert.Equal(t, createReq.Title, book.Title)
	assert.Equal(t, createReq.Author, book.Author)
	assert.Equal(t, createReq.Pages, book.Pages)

	// Verify book is in store
	storedBook, ok := bookStore.GetBook(book.ID)
	assert.True(t, ok)
	assert.Equal(t, book.Title, storedBook.Title)
}

func TestCreateBookInvalidInput(t *testing.T) {
	router := chi.NewRouter()
	bookStore := NewBookStore()
	router.Post("/books", CreateBookHandler(bookStore))

	// Test case 1: Missing title
	invalidReq := CreateBookRequest{Author: "Author", Pages: 10}
	body, _ := json.Marshal(invalidReq)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test case 2: Title too short
	invalidReq = CreateBookRequest{Title: "ab", Author: "Author", Pages: 10}
	body, _ = json.Marshal(invalidReq)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Test case 3: Pages less than 1
	invalidReq = CreateBookRequest{Title: "Valid Title", Author: "Author", Pages: 0}
	body, _ = json.Marshal(invalidReq)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	router := chi.NewRouter()
	bookStore := NewBookStore()
	router.Get("/books/{id}", GetBookHandler(bookStore))

	book1 := Book{ID: uuid.New().String(), Title: "Test Book 1", Author: "Author 1", Pages: 100}
	bookStore.CreateBook(book1)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", book1.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var book Book
	err := json.NewDecoder(rr.Body).Decode(&book)
	assert.NoError(t, err)
	assert.Equal(t, book1.ID, book.ID)
	assert.Equal(t, book1.Title, book.Title)
}

func TestGetBookNotFound(t *testing.T) {
	router := chi.NewRouter()
	bookStore := NewBookStore()
	router.Get("/books/{id}", GetBookHandler(bookStore))

	req, _ := http.NewRequest("GET", "/books/nonexistent-id", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
