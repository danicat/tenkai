package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	// Initialize router and clear book store
	r := chi.NewRouter()
	r.Get("/books", listBooksHandler)
	
	mu.Lock()
	bookStore = make(map[string]Book)
	mu.Unlock()

	// Test empty list
	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "[]\n", rr.Body.String())

	// Add a book and test again
	mu.Lock()
	bookStore["1"] = Book{
		ID:     "1",
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}
	mu.Unlock()

	req, _ = http.NewRequest("GET", "/books", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	expectedBody, _ := json.Marshal([]Book{bookStore["1"]})
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, string(expectedBody), rr.Body.String())
}

func TestCreateBook(t *testing.T) {
	// Initialize router and clear book store
	r := chi.NewRouter()
	r.Post("/books", createBookHandler)
	
	mu.Lock()
	bookStore = make(map[string]Book)
	mu.Unlock()

	// Test successful creation
	newBook := CreateBookRequest{
		Title:  "New Book Title",
		Author: "New Author",
		Pages:  200,
	}
	jsonBody, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdBook Book
	json.NewDecoder(rr.Body).Decode(&createdBook)
	assert.NotNil(t, createdBook.ID)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)
	assert.Equal(t, newBook.Pages, createdBook.Pages)

	mu.RLock()
	_, ok := bookStore[createdBook.ID]
	mu.RUnlock()
	assert.True(t, ok)

	// Test invalid input (missing title)
	invalidBook := CreateBookRequest{
		Author: "Invalid Author",
		Pages:  50,
	}
	jsonBody, _ = json.Marshal(invalidBook)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.True(t, strings.Contains(rr.Body.String(), "Key: 'CreateBookRequest.Title'") && strings.Contains(rr.Body.String(), "Error:Field validation for 'Title' failed on the 'required' tag"))

	// Test invalid input (pages < 1)
	invalidBook2 := CreateBookRequest{
		Title:  "Another Book",
		Author: "Another Author",
		Pages:  0,
	}
	jsonBody, _ = json.Marshal(invalidBook2)
	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.True(t, strings.Contains(rr.Body.String(), "Key: 'CreateBookRequest.Pages'") && strings.Contains(rr.Body.String(), "Error:Field validation for 'Pages' failed on the 'required' tag"))
}

func TestGetBook(t *testing.T) {
	// Initialize router and clear book store
	r := chi.NewRouter()
	r.Get("/books/{id}", getBookHandler)
	
	mu.Lock()
	bookStore = make(map[string]Book)
	bookStore["existing-id"] = Book{
		ID:     "existing-id",
		Title:  "Existing Book",
		Author: "Existing Author",
		Pages:  300,
	}
	mu.Unlock()

	// Test retrieving existing book
	req, _ := http.NewRequest("GET", "/books/existing-id", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	expectedBook, _ := json.Marshal(bookStore["existing-id"])
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, string(expectedBook), rr.Body.String())

	// Test non-existent book
	req, _ = http.NewRequest("GET", "/books/non-existent-id", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "Book not found\n", rr.Body.String())
}
