package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/models"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
	"github.com/Tarunika21/ImpactGraph/services/api/internal/services"
)

type AnalysisHandler struct {
	dependencyRepository *repositories.DependencyRepository
	reportRepository     *repositories.ReportRepository
}

type AnalysisRequest struct {
	ServiceName string             `json:"serviceName"`
	OldVersion  string             `json:"oldVersion"`
	NewVersion  string             `json:"newVersion"`
	OldSpec     models.OpenAPISpec `json:"oldSpec"`
	NewSpec     models.OpenAPISpec `json:"newSpec"`
}

type AnalysisResponse struct {
	ReportID   int64                     `json:"reportId"`
	Comparison models.ComparisonResponse `json:"comparison"`
	Impact     models.ImpactResponse     `json:"impact"`
}

func NewAnalysisHandler(
	dependencyRepository *repositories.DependencyRepository,
	reportRepository *repositories.ReportRepository,
) *AnalysisHandler {
	return &AnalysisHandler{
		dependencyRepository: dependencyRepository,
		reportRepository:     reportRepository,
	}
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
	request.OldVersion = strings.TrimSpace(request.OldVersion)
	request.NewVersion = strings.TrimSpace(request.NewVersion)

	if request.ServiceName == "" ||
		request.OldVersion == "" ||
		request.NewVersion == "" ||
		request.OldSpec.Paths == nil ||
		request.NewSpec.Paths == nil {
		http.Error(w, "serviceName, oldVersion, newVersion, oldSpec.paths and newSpec.paths are required", http.StatusBadRequest)
		return
	}

	if len(request.ServiceName) > 150 ||
		len(request.OldVersion) > 100 ||
		len(request.NewVersion) > 100 {
		http.Error(w, "Service name or version exceeds the allowed length", http.StatusBadRequest)
		return
	}

	comparison, err := services.CompareSpecsChecked(request.OldSpec, request.NewSpec)
	if err != nil {
		http.Error(w, "Contract comparison failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	consumers, err := h.dependencyRepository.GetImpactConsumers(
		r.Context(),
		request.ServiceName,
	)
	if err != nil {
		log.Printf("Failed to retrieve consumer dependencies: %v", err)
		http.Error(w, "Failed to retrieve consumer dependencies", http.StatusInternalServerError)
		return
	}

	impact := services.AnalyzeImpact(comparison.Changes, consumers)

	response := AnalysisResponse{
		Comparison: comparison,
		Impact:     impact,
	}

	resultJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to serialize analysis report: %v", err)
		http.Error(w, "Failed to serialize analysis report", http.StatusInternalServerError)
		return
	}

	report := models.ImpactReport{
		ServiceName: request.ServiceName,
		OldVersion:  request.OldVersion,
		NewVersion:  request.NewVersion,
		Result:      resultJSON,
	}

	savedReport, err := h.reportRepository.Create(r.Context(), report)
	if err != nil {
		log.Printf("Failed to save analysis report: %v", err)
		http.Error(w, "Failed to save analysis report", http.StatusInternalServerError)
		return
	}

	response.ReportID = savedReport.ID

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write analysis response: %v", err)
	}
}
