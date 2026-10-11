package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/services"
)

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.ComparisonRequest

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid comparison request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if request.OldSpec.Paths == nil || request.NewSpec.Paths == nil {
		http.Error(w, "Both specifications must contain paths", http.StatusBadRequest)
		return
	}

	result := services.CompareSpecs(request.OldSpec, request.NewSpec)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
