package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestListBooks(t *testing.T) {
	// Clear books before testing
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	// Add a dummy book
	dummyBook := Book{
		ID:    "123",
		Title:  "Test Book",
		Author: "Test Author",
		Pages: 100,
	}
	mu.Lock()
	books[dummyBook.ID] = dummyBook
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books", listBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var bookList []Book
	err := json.NewDecoder(rec.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 1)
	assert.Equal(t, dummyBook.Title, bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear books before testing
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Test successful creation
	createReq := CreateBookRequest{
		Title:  "New Book Title",
		Author: "New Author",
		Pages: 150,
	}
	reqBody, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var createdBook Book
	err := json.NewDecoder(rec.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.NotNil(t, createdBook.ID)
	assert.Equal(t, createReq.Title, createdBook.Title)
	assert.Equal(t, createReq.Author, createdBook.Author)
	assert.Equal(t, createReq.Pages, createdBook.Pages)

	// Test validation error (missing title)
	invalidReq := CreateBookRequest{
		Author: "Invalid Author",
		Pages: 10,
	}
	reqBody, _ = json.Marshal(invalidReq)

	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errs map[string]string
	err = json.NewDecoder(rec.Body).Decode(&errs)
	assert.NoError(t, err)
	assert.Contains(t, errs, "Title")
	assert.Contains(t, errs["Title"], "Title validation failed on the 'required' tag")

	// Test validation error (pages less than 1)
	invalidReq = CreateBookRequest{
		Title:  "Valid Title",
		Author: "Valid Author",
		Pages: 0,
	}
	reqBody, _ = json.Marshal(invalidReq)

	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	err = json.NewDecoder(rec.Body).Decode(&errs)
	assert.NoError(t, err)
	assert.Contains(t, errs, "Pages")
	assert.Contains(t, errs["Pages"], "Pages validation failed on the 'min' tag")

	// Test validation error (title too short)
	invalidReq = CreateBookRequest{
		Title:  "ab",
		Author: "Valid Author",
		Pages: 10,
	}
	reqBody, _ = json.Marshal(invalidReq)

	req, _ = http.NewRequest("POST", "/books", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	err = json.NewDecoder(rec.Body).Decode(&errs)
	assert.NoError(t, err)
	assert.Contains(t, errs, "Title")
	assert.Contains(t, errs["Title"], "Title validation failed on the 'min' tag")
}

func TestGetBook(t *testing.T) {
	// Clear books before testing
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Test book found
	dummyBook := Book{
		ID:    "found-123",
		Title:  "Found Book",
		Author: "Found Author",
		Pages: 200,
	}
	mu.Lock()
	books[dummyBook.ID] = dummyBook
	mu.Unlock()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", dummyBook.ID), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var retrievedBook Book
	err := json.NewDecoder(rec.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, dummyBook.Title, retrievedBook.Title)

	// Test book not found
	req, _ = http.NewRequest("GET", "/books/non-existent", nil)
	rec = httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Book not found")
}
