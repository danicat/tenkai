package bookstore

import "github.com/google/uuid"

// Book represents a book in the bookstore.
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

// CreateBookRequest represents the request body for creating a new book.
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// NewBook creates a new Book with a generated ID.
func NewBook(req CreateBookRequest) Book {
	return Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}
}
