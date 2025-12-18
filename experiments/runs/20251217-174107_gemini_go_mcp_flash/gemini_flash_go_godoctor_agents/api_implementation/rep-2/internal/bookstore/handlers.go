package bookstore

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// API represents the API service with its dependencies.
type API struct {
	store    *Store
	validator *validator.Validate
}

// NewAPI creates and returns a new API instance.
func NewAPI(store *Store) *API {
	return &API{
		store:    store,
		validator: validator.New(),
	}
}

// handleListBooks godoc
// @Summary List all books
// @Produce json
// @Success 200 {array} Book
// @Router /books [get]
func (api *API) handleListBooks(w http.ResponseWriter, r *http.Request) {
	books := api.store.ListBooks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// handleCreateBook godoc
// @Summary Create a book
// @Accept json
// @Produce json
// @Param book body CreateBookRequest true "Book to create"
// @Success 201 {object} Book
// @Failure 400 {string} string "Invalid input"
// @Router /books [post]
func (api *API) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := api.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := NewBook(req)
	api.store.AddBook(book)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// handleGetBook godoc
// @Summary Get a book by ID
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} Book
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [get]
func (api *API) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := api.store.GetBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
