package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"experiment/internal/book"
	"experiment/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Handler holds the dependencies for our HTTP handlers.
type Handler struct {
	BookStore *store.MemBookStore
	Validator *validator.Validate
}

// NewHandler creates a new Handler instance.
func NewHandler(bookStore *store.MemBookStore) *Handler {
	return &Handler{
		BookStore: bookStore,
		Validator: validator.New(),
	}
}

// ListBooks returns a list of all books.
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := h.BookStore.ListBooks()

	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateBook creates a new book.
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req book.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		// This handles validation errors from the 'validate' tags in the struct
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	bookID := uuid.New().String()
	newBook := book.Book{
		ID:     bookID,
		Title:  req.Title,
		Author: req.Author,
		Pages:  req.Pages,
	}

	if err := h.BookStore.CreateBook(newBook); err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newBook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetBook returns a single book by ID.
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	b, err := h.BookStore.GetBook(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
