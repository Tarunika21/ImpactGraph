package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
)

type ConsumerHandler struct {
	repository *repositories.ConsumerRepository
}

func NewConsumerHandler(repository *repositories.ConsumerRepository) *ConsumerHandler {
	return &ConsumerHandler{repository: repository}
}

func (h *ConsumerHandler) HandleConsumers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createConsumer(w, r)
	case http.MethodGet:
		h.getConsumers(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ConsumerHandler) createConsumer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request models.CreateConsumerRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	request.ID = strings.TrimSpace(request.ID)
	request.Name = strings.TrimSpace(request.Name)

	if request.ID == "" || request.Name == "" {
		http.Error(w, "Consumer ID and name are required", http.StatusBadRequest)
		return
	}

	consumer, err := h.repository.Create(r.Context(), request)
	if err != nil {
		if errors.Is(err, repositories.ErrConsumerExists) {
			http.Error(w, "Consumer already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create consumer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(consumer)
}

func (h *ConsumerHandler) getConsumers(w http.ResponseWriter, r *http.Request) {
	consumers, err := h.repository.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve consumers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(consumers)
}
