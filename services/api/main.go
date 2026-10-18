package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/database"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/handlers"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
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
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx := context.Background()
	dbPool, err := database.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()

	log.Println("Connected to PostgreSQL successfully")

	contractRepository := repositories.NewContractRepository(dbPool)
	contractHandler := handlers.NewContractHandler(contractRepository)

	consumerRepository := repositories.NewConsumerRepository(dbPool)
	consumerHandler := handlers.NewConsumerHandler(consumerRepository)

	dependencyRepository := repositories.NewDependencyRepository(dbPool)
	dependencyHandler := handlers.NewDependencyHandler(dependencyRepository)

	storedImpactHandler := handlers.NewStoredImpactHandler(dependencyRepository)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/compare", handlers.CompareHandler)
	mux.HandleFunc("/api/v1/impact", handlers.ImpactHandler)
	mux.HandleFunc("/api/v1/contracts", contractHandler.HandleContracts)
	mux.HandleFunc("/api/v1/consumers", consumerHandler.HandleConsumers)
	mux.HandleFunc("/api/v1/dependencies", dependencyHandler.HandleDependencies)
	mux.HandleFunc("/api/v1/impact/stored", storedImpactHandler.Handle)

	server := &http.Server{
		Addr:              "127.0.0.1:8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("ImpactGraph API running on http://localhost:8081")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
