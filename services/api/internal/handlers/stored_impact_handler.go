package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/services"
)

type StoredImpactHandler struct {
	repository *repositories.DependencyRepository
}

type StoredImpactRequest struct {
	ServiceName string          `json:"serviceName"`
	Changes     []models.Change `json:"changes"`
}

func NewStoredImpactHandler(repository *repositories.DependencyRepository) *StoredImpactHandler {
	return &StoredImpactHandler{repository: repository}
}

func (h *StoredImpactHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	defer r.Body.Close()

	var request StoredImpactRequest
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

	request.ServiceName = strings.TrimSpace(request.ServiceName)

	if request.ServiceName == "" || len(request.ServiceName) > 150 || request.Changes == nil {
		http.Error(w, "Service name and changes are required", http.StatusBadRequest)
		return
	}

	for _, change := range request.Changes {
		if strings.TrimSpace(change.Type) == "" ||
			strings.TrimSpace(change.Endpoint) == "" {
			http.Error(w, "Each change requires a type and endpoint", http.StatusBadRequest)
			return
		}
	}

	consumers, err := h.repository.GetImpactConsumers(
		r.Context(),
		request.ServiceName,
	)
	if err != nil {
		http.Error(w, "Failed to retrieve consumer dependencies", http.StatusInternalServerError)
		return
	}

	result := services.AnalyzeImpact(request.Changes, consumers)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
