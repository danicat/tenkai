package bookstore

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("book not found")
)

// Book represents a book in the bookstore.
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"min=1"`
}

// CreateBookRequest represents the request body for creating a new book.
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"min=1"`
}

// Service defines the interface for bookstore operations.
type Service interface {
	ListBooks(ctx context.Context) ([]Book, error)
	CreateBook(ctx context.Context, req CreateBookRequest) (Book, error)
	GetBook(ctx context.Context, id string) (Book, error)
}

// inMemoryService is an in-memory implementation of the Service interface.
type inMemoryService struct {
	mu    sync.RWMutex
	books map[string]Book
}

// NewService creates and returns a new in-memory bookstore service.
func NewService() Service {
	return &inMemoryService{
		books: make(map[string]Book),
	}
}

func (s *inMemoryService) ListBooks(ctx context.Context) ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allBooks []Book
	for _, book := range s.books {
		allBooks = append(allBooks, book)
	}

	return allBooks, nil
}

func (s *inMemoryService) CreateBook(ctx context.Context, req CreateBookRequest) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	book := Book{
		ID:     id,
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	s.books[id] = book

	return book, nil
}

func (s *inMemoryService) GetBook(ctx context.Context, id string) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book with ID %s: %w", id, ErrNotFound)
	}

	return book, nil
}
