package httpapi

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/excel"
	"homemaker-followup/internal/followup"
)

func TestHTTPWorkflow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "http.xlsx")
	workbook, err := excel.NewWorkbook(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := workbook.Save(); err != nil {
		t.Fatal(err)
	}
	workbook.Close()
	service := followup.NewService(path, nil)
	if _, err := service.Start(); err != nil {
		t.Fatal(err)
	}
	server := NewServer(service, nil, domain.ReminderSetting{ID: "default", Enabled: true}, time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status %d", recorder.Code)
	}
}
