package bookstore_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/bookstore"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func TestBookstoreAPI(t *testing.T) {
	store := bookstore.NewInMemBookStore()
	validate := validator.New()
	h := &bookstore.Handlers{
		Store:    store,
		Validate: validate,
	}

	r := chi.NewRouter()
	r.Get("/books", h.ListBooks)
	r.Post("/books", h.CreateBook)
	r.Get("/books/{id}", h.GetBook)

	var err error // Declare err once here

	t.Run("ListBooks_Empty", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, rec.Code)
		}

		expected := "[]"
		if rec.Body.String() != expected {
			t.Errorf("expected body %q; got %q", expected, rec.Body.String())
		}
	})

	t.Run("CreateBook_Success", func(t *testing.T) {
		newBook := bookstore.CreateBookRequest{
			Title:  "The Go Programming Language",
			Author: "Alan A. A. Donovan",
			Pages:  400,
		}
		body, _ := json.Marshal(newBook)
		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status %d; got %d", http.StatusCreated, rec.Code)
		}

		var createdBook bookstore.Book
		err = json.Unmarshal(rec.Body.Bytes(), &createdBook)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if createdBook.Title != newBook.Title || createdBook.Author != newBook.Author || createdBook.Pages != newBook.Pages || createdBook.ID == "" {
			t.Errorf("created book mismatch: %+v", createdBook)
		}
	})

	t.Run("CreateBook_InvalidInput", func(t *testing.T) {
		invalidBook := bookstore.CreateBookRequest{
			Title:  "Go", // Too short
			Author: "",   // Missing
			Pages:  0,    // Too low
		}
		body, _ := json.Marshal(invalidBook)
		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d; got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("GetBook_Success", func(t *testing.T) {
		// First, create a book to retrieve
		newBook := bookstore.CreateBookRequest{
			Title:  "Effective Go",
			Author: "The Go Authors",
			Pages:  100,
		}
		created, _ := store.CreateBook(context.Background(), newBook)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", created.ID), nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, rec.Code)
		}

		var retrievedBook bookstore.Book
		err = json.Unmarshal(rec.Body.Bytes(), &retrievedBook)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if retrievedBook.ID != created.ID || retrievedBook.Title != created.Title {
			t.Errorf("retrieved book mismatch: expected %+v, got %+v", created, retrievedBook)
		}
	})

	t.Run("GetBook_NotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/books/non-existent-id", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d; got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("ListBooks_AfterCreate", func(t *testing.T) {
		// Add another book to ensure ListBooks works with multiple items
		_, err = store.CreateBook(context.Background(), bookstore.CreateBookRequest{
			Title:  "Another Book",
			Author: "Another Author",
			Pages:  250,
		})
		if err != nil {
			t.Fatalf("failed to create book for ListBooks_AfterCreate: %v", err)
		}

		req, _ := http.NewRequest("GET", "/books", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, rec.Code)
		}

		var books []bookstore.Book
		err = json.Unmarshal(rec.Body.Bytes(), &books)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(books) != 3 {
			t.Errorf("expected 3 books; got %d", len(books))
		}
	})
}
