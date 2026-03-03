package main

import (
	"dian-downloader/internal/api"
	"log"
	"net/http"
)

func main() {
	h := api.NewHandler()

	http.HandleFunc("/api/v1/download", h.Download)

	log.Println("Server running on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
