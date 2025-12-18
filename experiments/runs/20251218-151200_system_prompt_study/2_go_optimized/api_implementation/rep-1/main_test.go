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

func TestListBooksEmpty(t *testing.T) {
	// Clear books before test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/books", ListBooks)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "[]\n", rr.Body.String())
}

func TestCreateBookValid(t *testing.T) {
	// Clear books before test
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	bookReq := CreateBookRequest{
		Title:  "The Hitchhiker's Guide to the Galaxy",
		Author: "Douglas Adams",
		Pages:  193,
	}
	body, _ := json.Marshal(bookReq)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/books", CreateBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var book Book
	err := json.Unmarshal(rr.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.NotEmpty(t, book.ID)
	assert.Equal(t, bookReq.Title, book.Title)
	assert.Equal(t, bookReq.Author, book.Author)
	assert.Equal(t, bookReq.Pages, book.Pages)

	// Verify the book is in storage
	mu.RLock()
	defer mu.RUnlock()
	_, ok := books[book.ID]
	assert.True(t, ok)
}

func TestCreateBookInvalidInput(t *testing.T) {
	testCases := []struct {
		name string
		requestBody CreateBookRequest
	}{
		{
			name: "Missing Title",
			requestBody: CreateBookRequest{
				Author: "Douglas Adams",
				Pages:  193,
			},
		},
		{
			name: "Short Title",
			requestBody: CreateBookRequest{
				Title:  "Hi",
				Author: "Douglas Adams",
				Pages:  193,
			},
		},
		{
			name: "Missing Author",
			requestBody: CreateBookRequest{
				Title:  "The Hitchhiker's Guide to the Galaxy",
				Pages:  193,
			},
		},
		{
			name: "Short Author",
			requestBody: CreateBookRequest{
				Title:  "The Hitchhiker's Guide to the Galaxy",
				Author: "DA",
				Pages:  193,
			},
		},
		{
			name: "Pages Less Than 1",
			requestBody: CreateBookRequest{
				Title:  "The Hitchhiker's Guide to the Galaxy",
				Author: "Douglas Adams",
				Pages:  0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/books", CreateBook)
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestGetBook(t *testing.T) {
	// Clear books and add a test book
	mu.Lock()
	books = make(map[string]Book)
	testBook := Book{
		ID:    "123",
		Title: "Test Book",
		Author: "Test Author",
		Pages: 100,
	}
	books[testBook.ID] = testBook
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/123", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var book Book
	err := json.Unmarshal(rr.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, testBook.ID, book.ID)
	assert.Equal(t, testBook.Title, book.Title)
}

func TestGetBookNotFound(t *testing.T) {
	// Clear books
	mu.Lock()
	books = make(map[string]Book)
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books/nonexistent", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestListBooksWithData(t *testing.T) {
	// Clear books and add a test book
	mu.Lock()
	books = make(map[string]Book)
	testBook1 := Book{ID: "1", Title: "Book One", Author: "Author One", Pages: 100}
	testBook2 := Book{ID: "2", Title: "Book Two", Author: "Author Two", Pages: 200}
	books[testBook1.ID] = testBook1
	books[testBook2.ID] = testBook2
	mu.Unlock()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/books", ListBooks)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var bookList []Book
	err := json.Unmarshal(rr.Body.Bytes(), &bookList)
	assert.NoError(t, err)
	assert.Len(t, bookList, 2)
}
