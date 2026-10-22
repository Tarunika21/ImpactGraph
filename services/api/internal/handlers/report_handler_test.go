package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleReportsRejectsPost(t *testing.T) {
	handler := &ReportHandler{}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/reports",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleReports(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status 405, got %d",
			recorder.Code,
		)
	}

	if recorder.Header().Get("Allow") != http.MethodGet {
		t.Errorf("expected Allow header to be GET")
	}
}

func TestHandleReportByIDRejectsInvalidID(t *testing.T) {
	handler := &ReportHandler{}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/reports/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleReportByID(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestHandleReportByIDRejectsPost(t *testing.T) {
	handler := &ReportHandler{}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/reports/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleReportByID(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status 405, got %d",
			recorder.Code,
		)
	}
}
