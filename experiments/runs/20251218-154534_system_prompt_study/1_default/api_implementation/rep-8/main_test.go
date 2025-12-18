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
	r := chi.NewRouter()
	r.Get("/books", ListBooksHandler(bookStore))

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var books []Book
	json.NewDecoder(rr.Body).Decode(&books)
	assert.Empty(t, books)

	// Add a book and test again
	bookStore.CreateBook(Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100})
	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.NewDecoder(rr.Body).Decode(&books)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)
}

func TestCreateBook(t *testing.T) {
	bookStore := NewBookStore()
	r := chi.NewRouter()
	r.Post("/books", CreateBookHandler(bookStore))

	// Valid request
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var book Book
	json.NewDecoder(rr.Body).Decode(&book)
	assert.NotEmpty(t, book.ID)
	assert.Equal(t, newBook.Title, book.Title)
	assert.Equal(t, newBook.Author, book.Author)
	assert.Equal(t, newBook.Pages, book.Pages)

	// Invalid request - missing title
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Invalid request - pages too low
	invalidBook2 := CreateBookRequest{Title: "Book", Author: "Author", Pages: 0}
	body, _ = json.Marshal(invalidBook2)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetBook(t *testing.T) {
	bookStore := NewBookStore()
	r := chi.NewRouter()
	r.Get("/books/{id}", GetBookHandler(bookStore))

	// Book not found
	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Create a book and retrieve it
	book := Book{ID: "abc-123", Title: "Found Book", Author: "Found Author", Pages: 300}
	bookStore.CreateBook(book)

	req, _ = http.NewRequest("GET", "/books/abc-123", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var retrievedBook Book
	json.NewDecoder(rr.Body).Decode(&retrievedBook)
	assert.Equal(t, book.ID, retrievedBook.ID)
	assert.Equal(t, book.Title, retrievedBook.Title)
}
