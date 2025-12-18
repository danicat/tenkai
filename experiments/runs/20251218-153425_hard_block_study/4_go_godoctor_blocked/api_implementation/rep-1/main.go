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
	var bookList []Book
	for _, book := range bs.books {
		bookList = append(bookList, book)
	}
	return bookList
}

func (bs *BookStore) GetBook(id string) (Book, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	book, ok := bs.books[id]
	return book, ok
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
	if err := http.ListenAndServe(":8080", r); err != nil && err != http.ErrServerClosed {
		fmt.Printf("server error: %s\n", err)
	}
}

func ListBooksHandler(bs *BookStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books := bs.ListBooks()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			fmt.Printf("error encoding response in ListBooksHandler: %s\n", err)
		}
	}
}

func CreateBookHandler(bs *BookStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateBookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		if err := json.NewEncoder(w).Encode(book); err != nil {
			fmt.Printf("error encoding response in CreateBookHandler: %s\n", err)
		}
	}
}

func GetBookHandler(bs *BookStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		book, ok := bs.GetBook(id)
		if !ok {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(book); err != nil {
			fmt.Printf("error encoding response in GetBookHandler: %s\n", err)
		}
	}
}
