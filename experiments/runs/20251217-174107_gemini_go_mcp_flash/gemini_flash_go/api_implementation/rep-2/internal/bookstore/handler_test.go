package bookstore_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/petruzalek/tiny-bookstore/internal/bookstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookstoreHandlers(t *testing.T) {
	logger := log.New(os.Stdout, "test-bookstore: ", log.LstdFlags)
	service := bookstore.NewService()
	validate := validator.New()

	handler := bookstore.NewHandler(service, validate, logger)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// Use httptest.NewServer to get a real HTTP server for integration tests
	server := httptest.NewServer(r)
	defer server.Close()

	client := server.Client()

	ctx := context.Background()

	// Test 1: ListBooks (initially empty)
	t.Run("ListBooks_Empty", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/books", nil)
		require.NoError(t, err)

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		var books []bookstore.Book
		err = json.NewDecoder(res.Body).Decode(&books)
		require.NoError(t, err)
		assert.Empty(t, books)
	})

	// Test 2: CreateBook (valid)
	var createdBook bookstore.Book
	t.Run("CreateBook_Valid", func(t *testing.T) {
		createReq := bookstore.CreateBookRequest{
			Title:  "The Go Programming Language",
			Author: "Alan A. A. Donovan",
			Pages:  400,
		}
		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/books", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		err = json.NewDecoder(res.Body).Decode(&createdBook)
		require.NoError(t, err)
		assert.NotEmpty(t, createdBook.ID)
		assert.Equal(t, createReq.Title, createdBook.Title)
		assert.Equal(t, createReq.Author, createdBook.Author)
		assert.Equal(t, createReq.Pages, createdBook.Pages)
	})

	// Test 3: ListBooks (with one book)
	t.Run("ListBooks_OneBook", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/books", nil)
		require.NoError(t, err)

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		var books []bookstore.Book
		err = json.NewDecoder(res.Body).Decode(&books)
		require.NoError(t, err)
		assert.Len(t, books, 1)
		assert.Equal(t, createdBook, books[0])
	})

	// Test 4: GetBook (found)
	t.Run("GetBook_Found", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/books/"+createdBook.ID, nil)
		require.NoError(t, err)

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		var book bookstore.Book
		err = json.NewDecoder(res.Body).Decode(&book)
		require.NoError(t, err)
		assert.Equal(t, createdBook, book)
	})

	// Test 5: GetBook (not found)
	t.Run("GetBook_NotFound", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/books/nonexistent-id", nil)
		require.NoError(t, err)

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	// Test 6: CreateBook (invalid input - missing title)
	t.Run("CreateBook_InvalidInput_MissingTitle", func(t *testing.T) {
		createReq := bookstore.CreateBookRequest{
			Author: "Invalid Author",
			Pages:  100,
		}
		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/books", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		expectedErrorMsg := "Key: 'CreateBookRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
		buf := new(strings.Builder)
		_, err = io.Copy(buf, res.Body)
		require.NoError(t, err)
		assert.True(t, strings.Contains(buf.String(), expectedErrorMsg), "Expected error message not found: %s in %s", expectedErrorMsg, buf.String())
	})

	// Test 7: CreateBook (invalid input - pages too low)
	t.Run("CreateBook_InvalidInput_PagesTooLow", func(t *testing.T) {
		createReq := bookstore.CreateBookRequest{
			Title:  "Valid Title",
			Author: "Valid Author",
			Pages:  0,
		}
		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/books", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		expectedErrorMsg := "Key: 'CreateBookRequest.Pages' Error:Field validation for 'Pages' failed on the 'min' tag"
		buf := new(strings.Builder)
		_, err = io.Copy(buf, res.Body)
		require.NoError(t, err)
		assert.True(t, strings.Contains(buf.String(), expectedErrorMsg), "Expected error message not found: %s in %s", expectedErrorMsg, buf.String())
	})
	
	// Test 8: CreateBook (invalid input - author too short)
	t.Run("CreateBook_InvalidInput_AuthorTooShort", func(t *testing.T) {
		createReq := bookstore.CreateBookRequest{
			Title:  "Valid Title",
			Author: "AB", // Less than 3 characters
			Pages:  10,
		}
		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/books", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		expectedErrorMsg := "Key: 'CreateBookRequest.Author' Error:Field validation for 'Author' failed on the 'min' tag"
		buf := new(strings.Builder)
		_, err = io.Copy(buf, res.Body)
		require.NoError(t, err)
		assert.True(t, strings.Contains(buf.String(), expectedErrorMsg), "Expected error message not found: %s in %s", expectedErrorMsg, buf.String())
	})
}
