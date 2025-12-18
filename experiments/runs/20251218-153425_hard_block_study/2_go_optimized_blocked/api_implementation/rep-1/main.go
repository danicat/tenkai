package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	bookStore := NewBookStore()
	h := &Handler{BookStore: bookStore}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", h.ListBooks)
	r.Post("/books", h.CreateBook)
	r.Get("/books/{id}", h.GetBook)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
