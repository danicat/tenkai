package bookstore

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler handles HTTP requests for the bookstore service.
type Handler struct {
	service  Service
	validate *validator.Validate
	logger   *log.Logger
}

// NewHandler creates and returns a new Handler.
func NewHandler(service Service, validate *validator.Validate, logger *log.Logger) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
		logger:   logger,
	}
}

// RegisterRoutes registers the bookstore API routes with the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/books", h.ListBooks)
	r.Post("/books", h.CreateBook)
	r.Get("/books/{id}", h.GetBook)
}

// ListBooks handles the GET /books request.
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.ListBooks(r.Context())
	if err != nil {
		h.logger.Printf("error listing books: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
		h.logger.Printf("error encoding books: %v", err)
		// Fallback to generic error, as we're already trying to write to the response writer
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// CreateBook handles the POST /books request.
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("error decoding create book request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Printf("validation error for create book request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.service.CreateBook(r.Context(), req)
	if err != nil {
		h.logger.Printf("error creating book: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(book); err != nil {
		h.logger.Printf("error encoding created book: %v", err)
		// Fallback to generic error, as we're already trying to write to the response writer
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// GetBook handles the GET /books/{id} request.
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	book, err := h.service.GetBook(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		h.logger.Printf("error getting book: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		h.logger.Printf("error encoding book: %v", err)
		// Fallback to generic error, as we're already trying to write to the response writer
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
