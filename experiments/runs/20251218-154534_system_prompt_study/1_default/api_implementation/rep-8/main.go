package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Book represents a book in the bookstore
type Book struct {
	ID    string `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// CreateBookRequest represents the request body for creating a new book
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	Author string `json:"author" validate:"required,min=3"`
	Pages  int    `json:"pages" validate:"required,min=1"`
}

// In-memory store for books
type BookStore struct {
	mu    sync.RWMutex
	books map[string]Book
}

func NewBookStore() *BookStore {
	return &BookStore{
		books: make(map[string]Book),
	}
}

func (bs *BookStore) ListBooks() []Book {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	var books []Book
	for _, book := range bs.books {
		books = append(books, book)
	}
	return books
}

func (bs *BookStore) GetBook(id string) (Book, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	book, found := bs.books[id]
	return book, found
}

func (bs *BookStore) CreateBook(book Book) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.books[book.ID] = book
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	bookStore := NewBookStore()

	r.Get("/books", ListBooksHandler(bookStore))
	r.Post("/books", CreateBookHandler(bookStore))
	r.Get("/books/{id}", GetBookHandler(bookStore))

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

func ListBooksHandler(bs *BookStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books := bs.ListBooks()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

func CreateBookHandler(bs *BookStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateBookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		book := Book{
			ID:     uuid.New().String(),
			Title:  req.Title,
			Author: req.Author,
			Pages:  req.Pages,
		}

		bs.CreateBook(book)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func GetBookHandler(bs *BookStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		book, found := bs.GetBook(id)
		if !found {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}
