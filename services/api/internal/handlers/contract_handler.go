package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
)

type ContractHandler struct {
	repository *repositories.ContractRepository
}

func NewContractHandler(repository *repositories.ContractRepository) *ContractHandler {
	return &ContractHandler{repository: repository}
}

func (h *ContractHandler) HandleContracts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createContract(w, r)
	case http.MethodGet:
		h.getContracts(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ContractHandler) createContract(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	defer r.Body.Close()

	var request models.CreateContractRequest
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
	request.Version = strings.TrimSpace(request.Version)

	if request.ServiceName == "" || request.Version == "" ||
		len(request.ServiceName) > 150 || len(request.Version) > 100 ||
		len(request.Specification) == 0 ||
		!json.Valid(request.Specification) {
		http.Error(w, "Invalid contract data", http.StatusBadRequest)
		return
	}

	var specification map[string]json.RawMessage
	if err := json.Unmarshal(request.Specification, &specification); err != nil || specification == nil {
		http.Error(w, "Specification must be a JSON object", http.StatusBadRequest)
		return
	}

	contract, err := h.repository.Create(r.Context(), request)
	if err != nil {
		if errors.Is(err, repositories.ErrContractExists) {
			http.Error(w, "Contract version already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to save contract", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(contract)
}

func (h *ContractHandler) getContracts(w http.ResponseWriter, r *http.Request) {
	contracts, err := h.repository.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve contracts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(contracts)
}
