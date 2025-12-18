package bookstore

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates and configures a new Chi router for the bookstore API.
func NewRouter(api *API) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", api.handleListBooks)
	r.Post("/books", api.handleCreateBook)
	r.Get("/books/{id}", api.handleGetBook)

	return r
}
