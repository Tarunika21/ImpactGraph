package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tarunika21/ImpactGraph/services/api/internal/repositories"
)

type ReportHandler struct {
	repository *repositories.ReportRepository
}

func NewReportHandler(repository *repositories.ReportRepository) *ReportHandler {
	return &ReportHandler{repository: repository}
}

func (h *ReportHandler) HandleReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reports, err := h.repository.GetAll(r.Context())
	if err != nil {
		log.Printf("Failed to retrieve reports: %v", err)
		http.Error(w, "Failed to retrieve reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reports); err != nil {
		log.Printf("Failed to write reports response: %v", err)
	}
}
