package bookstore

import (
	"fmt"
	"sync"
)

// Store provides an in-memory storage for books.
type Store struct {
	mu    sync.RWMutex
	books map[string]Book
}

// NewStore creates and returns a new in-memory Store.
func NewStore() *Store {
	return &Store{
		books: make(map[string]Book),
	}
}

// AddBook adds a new book to the store.
func (s *Store) AddBook(book Book) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.books[book.ID] = book
}

// GetBook returns a book by its ID. Returns an error if the book is not found.
func (s *Store) GetBook(id string) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book with ID %s not found", id)
	}
	return book, nil
}

// ListBooks returns all books in the store.
func (s *Store) ListBooks() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books
}
