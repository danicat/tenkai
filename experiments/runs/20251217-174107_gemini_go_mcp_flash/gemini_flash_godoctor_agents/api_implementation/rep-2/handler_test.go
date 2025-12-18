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
	// Clear books map for a clean test environment
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Get("/books", listBooks)

	req := httptest.NewRequest("GET", "/books", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var bookList []Book
	err := json.NewDecoder(rec.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Empty(t, bookList)

	// Add a book and retest
	book := Book{
		ID:     uuid.New().String(),
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}
	books[book.ID] = book

	req = httptest.NewRequest("GET", "/books", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.NewDecoder(rec.Body).Decode(&bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 1)
	assert.Equal(t, book.Title, bookList[0].Title)
}

func TestCreateBook(t *testing.T) {
	// Clear books map for a clean test environment
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Post("/books", createBook)

	// Valid request
	validReq := CreateBookRequest{
		Title:  "New Book",
		Author: "John Doe",
		Pages:  250,
	}
	body, _ := json.Marshal(validReq)
	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var createdBook Book
	err := json.NewDecoder(rec.Body).Decode(&createdBook)
	assert.NoError(t, err)
	assert.NotNil(t, createdBook.ID)
	assert.Equal(t, validReq.Title, createdBook.Title)
	assert.Equal(t, validReq.Author, createdBook.Author)
	assert.Equal(t, validReq.Pages, createdBook.Pages)

	_, err = uuid.Parse(createdBook.ID)
	assert.NoError(t, err, "Generated ID should be a valid UUID")

	// Invalid request - missing title
	invalidReq := CreateBookRequest{
		Author: "Jane Doe",
		Pages:  120,
	}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Contains(t, errResp.Message, "Title' Error:Field validation for 'Title' failed on the 'required' tag")

	// Invalid request - title too short
	invalidReq = CreateBookRequest{
		Title:  "ab",
		Author: "Jane Doe",
		Pages:  120,
	}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Contains(t, errResp.Message, "Title' Error:Field validation for 'Title' failed on the 'min' tag")

	// Invalid request - pages too small
	invalidReq = CreateBookRequest{
		Title:  "Valid Title",
		Author: "Valid Author",
		Pages:  0,
	}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Contains(t, errResp.Message, "Pages' Error:Field validation for 'Pages' failed on the 'required' tag")
}

func TestGetBook(t *testing.T) {
	// Clear books map for a clean test environment
	books = make(map[string]Book)

	r := chi.NewRouter()
	r.Get("/books/{id}", getBook)

	// Add a book for testing
	existingBook := Book{
		ID:     uuid.New().String(),
		Title:  "Existing Book",
		Author: "Author X",
		Pages:  300,
	}
	books[existingBook.ID] = existingBook

	// Get existing book
	req := httptest.NewRequest("GET", fmt.Sprintf("/books/%s", existingBook.ID), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var retrievedBook Book
	err := json.NewDecoder(rec.Body).Decode(&retrievedBook)
	assert.NoError(t, err)
	assert.Equal(t, existingBook.ID, retrievedBook.ID)
	assert.Equal(t, existingBook.Title, retrievedBook.Title)

	// Get non-existent book
	nonExistentID := uuid.New().String()
	req = httptest.NewRequest("GET", fmt.Sprintf("/books/%s", nonExistentID), nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Book not found")
}
