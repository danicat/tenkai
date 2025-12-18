package store

import (
	"fmt"
	"sync"

	"experiment/internal/book"
)

// BookStore defines the interface for book storage operations.
type BookStore interface {
	ListBooks() []book.Book
	GetBook(id string) (book.Book, error)
	CreateBook(book book.Book) error
}

// MemBookStore is an in-memory implementation of BookStore.
type MemBookStore struct {
	mu    sync.RWMutex
	books map[string]book.Book
}

// NewMemBookStore creates a new in-memory book store.
func NewMemBookStore() *MemBookStore {
	return &MemBookStore{
		books: make(map[string]book.Book),
	}
}

// ListBooks returns all books in the store.
func (s *MemBookStore) ListBooks() []book.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]book.Book, 0, len(s.books))
	for _, b := range s.books {
		books = append(books, b)
	}
	return books
}

// GetBook returns a book by its ID.
func (s *MemBookStore) GetBook(id string) (book.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.books[id]
	if !ok {
		return book.Book{}, fmt.Errorf("book with ID %s not found", id)
	}
	return b, nil
}

// CreateBook adds a new book to the store.
func (s *MemBookStore) CreateBook(b book.Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[b.ID]; ok {
		return fmt.Errorf("book with ID %s already exists", b.ID)
	}
	s.books[b.ID] = b
	return nil
}
