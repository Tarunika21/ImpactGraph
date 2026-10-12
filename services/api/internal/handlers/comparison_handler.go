package handlers

import (
	"encoding/json"
	"io"
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

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Request must contain exactly one JSON object", http.StatusBadRequest)
		return
	}

	if request.OldSpec.Paths == nil || request.NewSpec.Paths == nil {
		http.Error(w, "Both specifications must contain paths", http.StatusBadRequest)
		return
	}

	result, err := services.CompareSpecsChecked(request.OldSpec, request.NewSpec)
	if err != nil {
		http.Error(w, "Schema resolution failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}
