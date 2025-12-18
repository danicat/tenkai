package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Book represents a book in the bookstore.
type Book struct {
	ID     string `json:"id"`
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

// NewBook creates a new Book from a CreateBookRequest, generating a unique ID.
func NewBook(req CreateBookRequest) Book {
	return Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}
