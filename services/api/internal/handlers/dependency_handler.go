package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
)

type DependencyHandler struct {
	repository *repositories.DependencyRepository
}

func NewDependencyHandler(repository *repositories.DependencyRepository) *DependencyHandler {
	return &DependencyHandler{repository: repository}
}

func (h *DependencyHandler) HandleDependencies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createDependency(w, r)
	case http.MethodGet:
		h.getDependencies(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DependencyHandler) createDependency(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request models.CreateConsumerDependencyRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	request.ConsumerID = strings.TrimSpace(request.ConsumerID)
	request.ServiceName = strings.TrimSpace(request.ServiceName)
	request.Endpoint = strings.TrimSpace(request.Endpoint)

	if request.ConsumerID == "" || request.ServiceName == "" || request.Endpoint == "" {
		http.Error(w, "Consumer ID, service name and endpoint are required", http.StatusBadRequest)
		return
	}

	var fields []string
	if err := json.Unmarshal(request.Fields, &fields); err != nil || len(fields) == 0 {
		http.Error(w, "Fields must be a non-empty array of strings", http.StatusBadRequest)
		return
	}

	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			http.Error(w, "Field names cannot be empty", http.StatusBadRequest)
			return
		}
	}

	dependency, err := h.repository.Create(r.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrDependencyExists):
			http.Error(w, "Dependency already exists", http.StatusConflict)
		case errors.Is(err, repositories.ErrConsumerNotFound):
			http.Error(w, "Consumer not found", http.StatusNotFound)
		default:
			http.Error(w, "Failed to create dependency", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dependency)
}

func (h *DependencyHandler) getDependencies(w http.ResponseWriter, r *http.Request) {
	dependencies, err := h.repository.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve dependencies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dependencies)
}
