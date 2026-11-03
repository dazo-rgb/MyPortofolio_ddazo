package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dazo-rgb/ddazo-portfolio/internal/delivery/http"
	"github.com/dazo-rgb/ddazo-portfolio/internal/repository/memory"
	"github.com/dazo-rgb/ddazo-portfolio/internal/usecase"
)

func main() {
	messageRepo := memory.NewMessageRepository()
	contactService := usecase.NewContactService(messageRepo)
	handler := httpdelivery.NewHandler(contactService, "./web")

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	log.Printf("ddazo portfolio running at http://localhost:%s", addr)
	if err := http.ListenAndServe(":"+addr, handler); err != nil {
		log.Fatal(err)
	}
}
