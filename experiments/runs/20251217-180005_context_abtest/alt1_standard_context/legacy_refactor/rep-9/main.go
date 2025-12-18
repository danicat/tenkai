package main

import (
	"log"
	"net/http"

	"github.com/petruzalek/rep-9/internal/api"
	"github.com/petruzalek/rep-9/internal/service"
)

func main() {
	svc := service.NewService()
	handler := api.NewHandler(svc)
	handler.RegisterRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}