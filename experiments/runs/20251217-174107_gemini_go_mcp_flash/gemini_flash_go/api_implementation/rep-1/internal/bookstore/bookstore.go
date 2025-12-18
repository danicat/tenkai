package bookstore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Book represents a book in the store.
type Book struct {
	ID    string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// CreateBookRequest represents the request body for creating a new book.
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// Store provides methods for interacting with books.
type Store struct {
	books map[string]Book
	mu    sync.RWMutex
	validator *validator.Validate
}

// NewStore creates and returns a new Store.
func NewStore() *Store {
	return &Store{
		books: make(map[string]Book),
		validator: validator.New(),
	}
}

// ListBooks returns all books in the store.
func (s *Store) ListBooks(ctx context.Context) ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bookList []Book
	for _, book := range s.books {
		bookList = append(bookList, book)
	}
	return bookList, nil
}

// CreateBook adds a new book to the store.
func (s *Store) CreateBook(ctx context.Context, req CreateBookRequest) (Book, error) {
	if err := s.validator.Struct(req); err != nil {
		return Book{}, err // Return validation errors directly
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	book := Book{
		ID:    id,
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	s.books[id] = book
	return book, nil
}

// GetBook returns a book by its ID.
func (s *Store) GetBook(ctx context.Context, id string) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book not found")
	}
	return book, nil
}

// Server holds the HTTP handlers and store.
type Server struct {
	Store *Store
}

// NewServer creates and returns a new Server.
func NewServer() *Server {
	return &Server{
		Store: NewStore(),
	}
}

// Routes returns an http.Handler that serves the bookstore API.
func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/books", s.handleListBooks())
	r.Post("/books", s.handleCreateBook())
	r.Get("/books/{id}", s.handleGetBook())

	return r
}

func (s *Server) handleListBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := s.Store.ListBooks(r.Context())
		if err != nil {
			http.Error(w, "Failed to list books", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

func (s *Server) handleCreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateBookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		book, err := s.Store.CreateBook(r.Context(), req)
		if err != nil {
			if _, ok := err.(validator.ValidationErrors); ok {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "Failed to create book", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func (s *Server) handleGetBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Book ID is required", http.StatusBadRequest)
			return
		}

		book, err := s.Store.GetBook(r.Context(), id)
		if err != nil {
			if err.Error() == "book not found" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "Failed to get book", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}

