package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBookstoreAPI(t *testing.T) {
	srv := NewServer()
	r := srv.setupRoutes()
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	client := testServer.Client()

	// Test 1: GET /books (empty initially)
	t.Run("List books initially empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/books", nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		var books []Book
		if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(books) != 0 {
			t.Errorf("Expected 0 books, got %d", len(books))
		}
	})

	var createdBook Book

	// Test 2: POST /books (create a book)
	t.Run("Create a book successfully", func(t *testing.T) {
		createReq := CreateBookRequest{
			Title:  "The Go Programming Language",
			Author: "Alan A. A. Donovan",
			Pages:  500,
		}
		body, _ := json.Marshal(createReq)

		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/books", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, res.StatusCode)
		}

		if err := json.NewDecoder(res.Body).Decode(&createdBook); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if createdBook.ID == "" {
			t.Error("Expected book ID to be generated, got empty")
		}
		if createdBook.Title != createReq.Title || createdBook.Author != createReq.Author || createdBook.Pages != createReq.Pages {
			t.Errorf("Created book mismatch: got %+v, want %+v", createdBook, createReq)
		}
	})

	// Test 3: GET /books (check if the created book is present)
	t.Run("List books after creation", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/books", nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		var books []Book
		if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(books) != 1 {
			t.Fatalf("Expected 1 book, got %d", len(books))
		}
		if !reflect.DeepEqual(books[0], createdBook) {
			t.Errorf("Listed book mismatch: got %+v, want %+v", books[0], createdBook)
		}
	})

	// Test 4: GET /books/{id} (get existing book)
	t.Run("Get existing book by ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/books/"+createdBook.ID, nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		var book Book
		if err := json.NewDecoder(res.Body).Decode(&book); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !reflect.DeepEqual(book, createdBook) {
			t.Errorf("Retrieved book mismatch: got %+v, want %+v", book, createdBook)
		}
	})

	// Test 5: GET /books/{id} (book not found)
	t.Run("Get non-existent book by ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/books/non-existent-id", nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, res.StatusCode)
		}
	})

	// Test 6: POST /books (invalid input - title too short)
	t.Run("Create book with invalid title (too short)", func(t *testing.T) {
		createReq := CreateBookRequest{
			Title:  "Go", // Too short
			Author: "Someone",
			Pages:  100,
		}
		body, _ := json.Marshal(createReq)

		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/books", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

	// Test 7: POST /books (invalid input - pages less than 1)
	t.Run("Create book with invalid pages (less than 1)", func(t *testing.T) {
		createReq := CreateBookRequest{
			Title:  "A Valid Title",
			Author: "A Valid Author",
			Pages:  0, // Invalid
		}
		body, _ := json.Marshal(createReq)

		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/books", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

}
