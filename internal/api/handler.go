package api

import (
	"bytes"
	"dian-downloader/internal/client"
	"dian-downloader/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	dian *client.DianClient
}

func NewHandler() *Handler {
	return &Handler{dian: client.NewDianClient()}
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	var req models.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	go h.process(req)

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "processing"}); err != nil {
		log.Printf("Error sending JSON response to client: %v", err)
	}
}

func (h *Handler) process(req models.DownloadRequest) {
	start := time.Now()
	path := fmt.Sprintf("downloads/%s.pdf", req.DocumentKey)

	if err := os.MkdirAll("downloads", os.ModePerm); err != nil {
		log.Printf("Error creating directory: %v", err)
	}

	err := h.dian.DownloadPDF(req.DocumentKey, path)

	if req.WebhookURL != "" {
		payload := models.WebhookPayload{
			DocumentKey:   req.DocumentKey,
			Status:        "done",
			Success:       err == nil,
			FilePath:      path,
			ExecutionTime: time.Since(start).String(),
		}

		if err == nil {
			fileData, readErr := os.ReadFile(path) //
			if readErr != nil {
				log.Printf("Error reading file for webhook: %v", readErr)
				payload.Error = "Error reading file for transmission"
			} else {
				payload.FileContent = base64.StdEncoding.EncodeToString(fileData)
			}
		} else {
			payload.Error = err.Error()
		}

		data, _ := json.Marshal(payload)
		resp, err := http.Post(req.WebhookURL, "application/json", bytes.NewBuffer(data))

		if err != nil {
			payload.Error = err.Error()
			log.Printf("Error sending webhook: %v", err)
			return
		}

		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Printf("Error closing webhook response body: %v", closeErr)
			}
		}()

		log.Printf("Webhook enviado, respuesta: %s | Tiempo: %s", resp.Status, payload.ExecutionTime)
	}
}
