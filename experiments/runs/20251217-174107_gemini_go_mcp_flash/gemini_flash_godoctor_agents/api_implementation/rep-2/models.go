package main

import "github.com/go-playground/validator/v10"

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

// ErrorResponse represents a generic error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}
