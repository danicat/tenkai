package main

import (
	"sync"
)

type BookStore struct {
	mu    sync.RWMutex
	books map[string]Book
}

func NewBookStore() *BookStore {
	return &BookStore{
		books: make(map[string]Book),
	}
}

func (bs *BookStore) AddBook(book Book) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.books[book.ID] = book
}

func (bs *BookStore) GetBook(id string) (Book, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	book, ok := bs.books[id]
	return book, ok
}

func (bs *BookStore) GetAllBooks() []Book {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	books := make([]Book, 0, len(bs.books))
	for _, book := range bs.books {
		books = append(books, book)
	}
	return books
}
