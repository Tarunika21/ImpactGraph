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

type AnalysisHandler struct {
	repository *repositories.DependencyRepository
}

type AnalysisRequest struct {
	ServiceName string             `json:"serviceName"`
	OldSpec     models.OpenAPISpec `json:"oldSpec"`
	NewSpec     models.OpenAPISpec `json:"newSpec"`
}

type AnalysisResponse struct {
	Comparison models.ComparisonResponse `json:"comparison"`
	Impact     models.ImpactResponse     `json:"impact"`
}

func NewAnalysisHandler(repository *repositories.DependencyRepository) *AnalysisHandler {
	return &AnalysisHandler{repository: repository}
}

func (h *AnalysisHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	var request AnalysisRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid analysis request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Request must contain exactly one JSON object", http.StatusBadRequest)
		return
	}

	request.ServiceName = strings.TrimSpace(request.ServiceName)
	if request.ServiceName == "" || request.OldSpec.Paths == nil || request.NewSpec.Paths == nil {
		http.Error(w, "serviceName, oldSpec.paths and newSpec.paths are required", http.StatusBadRequest)
		return
	}

	comparison, err := services.CompareSpecsChecked(request.OldSpec, request.NewSpec)
	if err != nil {
		http.Error(w, "Contract comparison failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	consumers, err := h.repository.GetImpactConsumers(r.Context(), request.ServiceName)
	if err != nil {
		http.Error(w, "Failed to retrieve consumer dependencies", http.StatusInternalServerError)
		return
	}

	impact := services.AnalyzeImpact(comparison.Changes, consumers)

	response := AnalysisResponse{
		Comparison: comparison,
		Impact:     impact,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
