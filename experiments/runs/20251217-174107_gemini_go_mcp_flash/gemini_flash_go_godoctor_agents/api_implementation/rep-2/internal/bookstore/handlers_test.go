package bookstore

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestListBooks(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	// Add some dummy books
	store.AddBook(NewBook(CreateBookRequest{Title: "Book 1", Author: "Author 1", Pages: 100}))
	store.AddBook(NewBook(CreateBookRequest{Title: "Book 2", Author: "Author 2", Pages: 200}))

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	api.handleListBooks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var books []Book
	if err := json.NewDecoder(rr.Body).Decode(&books); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("handler returned wrong number of books: got %v want %v", len(books), 2)
	}
}

func TestCreateBook(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	tests := []struct {
		name       string
		reqBody    CreateBookRequest
		expectStatus int
		expectError  bool
	}{
		{
			name:       "Valid book",
			reqBody:    CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 150},
			expectStatus: http.StatusCreated,
			expectError:  false,
		},
		{
			name:       "Invalid title (too short)",
			reqBody:    CreateBookRequest{Title: "Nu", Author: "New Author", Pages: 150},
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:       "Invalid pages (zero)",
			reqBody:    CreateBookRequest{Title: "New Book", Author: "New Author", Pages: 0},
			expectStatus: http.StatusBadRequest,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			api.handleCreateBook(rr, req)

			if status := rr.Code; status != tt.expectStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectStatus)
			}

			if !tt.expectError {
				var book Book
				if err := json.NewDecoder(rr.Body).Decode(&book); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}
				if book.Title != tt.reqBody.Title {
					t.Errorf("handler returned unexpected book title: got %v want %v", book.Title, tt.reqBody.Title)
				}
				if book.ID == "" {
					t.Error("handler returned book with empty ID")
				}
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	store := NewStore()
	api := NewAPI(store)

	existingBook := NewBook(CreateBookRequest{Title: "Existing Book", Author: "Existing Author", Pages: 300})
	store.AddBook(existingBook)

	// Create a new router for testing with URL parameters
	r := chi.NewRouter()
	r.Get("/books/{id}", api.handleGetBook)

	tests := []struct {
		name       string
		bookID     string
		expectStatus int
		expectError  bool
	}{
		{
			name:       "Book found",
			bookID:     existingBook.ID,
			expectStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:       "Book not found",
			bookID:     "non-existent-id",
			expectStatus: http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/books/"+tt.bookID, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req) // Use the router to handle the request for URL params

			if status := rr.Code; status != tt.expectStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectStatus)
			}

			if !tt.expectError {
				var book Book
				if err := json.NewDecoder(rr.Body).Decode(&book); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}
				if book.ID != existingBook.ID {
					t.Errorf("handler returned unexpected book ID: got %v want %v", book.ID, existingBook.ID)
				}
			}
		})
	}
}
