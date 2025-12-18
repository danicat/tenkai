package bookstore

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// A Handlers holds all the http handlers for the book service.
type Handlers struct {
	Store    BookStore
	Validate *validator.Validate
}

// ListBooks returns a list of all books.
func (h *Handlers) ListBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	books, err := h.Store.ListBooks(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(books)
	if err != nil {
		log.Printf("error marshalling books: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(out); err != nil {
		log.Printf("error writing response in ListBooks: %v", err)
	}
}

// CreateBook creates a new book.
func (h *Handlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateBookRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.Store.CreateBook(ctx, req)
	if err != nil {
		log.Printf("error creating book: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return		
	}

	out, err := json.Marshal(book)
	if err != nil {
		log.Printf("error marshalling book: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(out); err != nil {
		log.Printf("error writing response in CreateBook: %v", err)
	}
}

// GetBook returns a single book by ID.
func (h *Handlers) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	book, err := h.Store.GetBook(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		log.Printf("error getting book: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(book)
	if err != nil {
		log.Printf("error marshalling book: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(out); err != nil {
		log.Printf("error writing response in GetBook: %v", err)
	}
}