package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	app := NewApp()
	r := app.routes()

	// Test empty list
	req, _ := http.NewRequest("GET", "/books", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))
	assert.Equal(t, "[]\n", rw.Body.String())

	// Add a book and test populated list
	app.mu.Lock()
	app.books["1"] = Book{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}
	app.mu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	expected := []Book{{ID: "1", Title: "Test Book", Author: "Test Author", Pages: 100}}
	var actual []Book
	err := json.NewDecoder(rw.Body).Decode(&actual)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, actual)
}

func TestCreateBook(t *testing.T) {
	app := NewApp()
	r := app.routes()

	// Test successful creation
	newBook := CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 200}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusCreated, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	var createdBook Book
	err := json.NewDecoder(rw.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)
	assert.Equal(t, newBook.Pages, createdBook.Pages)

	// Verify it's in the store
	app.mu.RLock()
	_, ok := app.books[createdBook.ID]
	app.mu.RUnlock()
	assert.True(t, ok)

	// Test validation error (missing title)
	invalidBook := CreateBookRequest{Author: "Invalid Author", Pages: 50}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	assert.Contains(t, rw.Body.String(), "Validation error")

	// Test validation error (pages too low)
	invalidBook = CreateBookRequest{Title: "Valid Title", Author: "Valid Author", Pages: 0}
	body, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
	assert.Contains(t, rw.Body.String(), "Validation error")
}

func TestGetBook(t *testing.T) {
	app := NewApp()
	r := app.routes()

	// Add a book
	existingBook := Book{ID: "test-123", Title: "Existing Book", Author: "Existing Author", Pages: 300}
	app.mu.Lock()
	app.books[existingBook.ID] = existingBook
	app.mu.Unlock()

	// Test found book
	req, _ := http.NewRequest("GET", "/books/test-123", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	var actual Book
	err := json.NewDecoder(rw.Body).Decode(&actual)
	assert.NoError(t, err)
	assert.Equal(t, existingBook, actual)

	// Test not found
	req, _ = http.NewRequest("GET", "/books/non-existent", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Code)
	assert.Contains(t, rw.Body.String(), "Book not found")
}
