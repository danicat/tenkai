package main

import (
	"log"
	"net/http"

	"legacy_refactor/rep-1/internal/api"
	"legacy_refactor/rep-1/internal/service"
)

func main() {
	svc := service.NewService()
	h := api.NewHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
