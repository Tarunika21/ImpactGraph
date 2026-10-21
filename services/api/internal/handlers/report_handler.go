package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	serviceName := strings.TrimSpace(r.URL.Query().Get("serviceName"))

	var (
		reportsErr error
	)

	if serviceName != "" {
		reports, err := h.repository.GetByService(r.Context(), serviceName)
		reportsErr = err
		if reportsErr != nil {
			log.Printf("Failed to retrieve reports: %v", reportsErr)
			http.Error(w, "Failed to retrieve reports", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(reports)
		return
	}

	reports, err := h.repository.GetAll(r.Context())
	if err != nil {
		log.Printf("Failed to retrieve reports: %v", err)
		http.Error(w, "Failed to retrieve reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) HandleReportByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/api/v1/reports/")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	report, err := h.repository.GetByID(r.Context(), id)
	if errors.Is(err, repositories.ErrReportNotFound) {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to retrieve report: %v", err)
		http.Error(w, "Failed to retrieve report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}
