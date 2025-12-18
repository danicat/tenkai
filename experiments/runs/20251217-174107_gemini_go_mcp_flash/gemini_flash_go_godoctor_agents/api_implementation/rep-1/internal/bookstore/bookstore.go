package bookstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var ErrBookNotFound = fmt.Errorf("book not found")

// Book represents a book in the bookstore.
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

// BookStore defines the interface for managing books.
type BookStore interface {
	ListBooks(ctx context.Context) ([]Book, error)
	CreateBook(ctx context.Context, req CreateBookRequest) (Book, error)
	GetBook(ctx context.Context, id string) (Book, error)
}

// InMemBookStore is an in-memory implementation of BookStore.
type InMemBookStore struct {
	mu    sync.RWMutex
	books map[string]Book
}

// NewInMemBookStore creates a new InMemBookStore.
func NewInMemBookStore() *InMemBookStore {
	return &InMemBookStore{
		books: make(map[string]Book),
	}
}

// ListBooks returns all books in the store.
func (s *InMemBookStore) ListBooks(ctx context.Context) ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bookList []Book
	if len(s.books) == 0 {
		return []Book{}, nil // Return an empty slice if no books
	}

	for _, book := range s.books {
		bookList = append(bookList, book)
	}
	return bookList, nil
}

// CreateBook creates a new book and adds it to the store.
func (s *InMemBookStore) CreateBook(ctx context.Context, req CreateBookRequest) (Book, error) {
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

// GetBook retrieves a book by its ID.
func (s *InMemBookStore) GetBook(ctx context.Context, id string) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}
	return book, nil
}
