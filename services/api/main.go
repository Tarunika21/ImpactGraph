package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/handlers"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Service:   "impactgraph-api",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write health response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/compare", handlers.CompareHandler)
	mux.HandleFunc("/api/v1/impact", handlers.ImpactHandler)

	server := &http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}

	log.Println("ImpactGraph API running on http://localhost:8081")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
