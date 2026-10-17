package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/services"
)

func ImpactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	defer r.Body.Close()

	var request models.ImpactRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Only one JSON object is allowed", http.StatusBadRequest)
		return
	}

	if request.Changes == nil || request.Consumers == nil {
		http.Error(w, "changes and consumers are required", http.StatusBadRequest)
		return
	}

	for _, consumer := range request.Consumers {
		if consumer.ID == "" || consumer.Endpoint == "" {
			http.Error(w, "Consumer ID and endpoint are required", http.StatusBadRequest)
			return
		}
	}

	result := services.AnalyzeImpact(request.Changes, request.Consumers)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
