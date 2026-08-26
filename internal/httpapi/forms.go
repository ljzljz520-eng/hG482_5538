package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"homemaker-followup/internal/domain"
)

type RecordForm struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id"`
	ClientName      string `json:"client_name"`
	ServiceTypeID   string `json:"service_type_id"`
	ServiceTypeName string `json:"service_type_name"`
	CaregiverID     string `json:"caregiver_id"`
	CaregiverName   string `json:"caregiver_name"`
	VisitDate       string `json:"visit_date"`
	NextFollowUp    string `json:"next_follow_up"`
	Score           int    `json:"score"`
	Comment         string `json:"comment"`
	Notes           string `json:"notes"`
}

func decodeForm(r *http.Request) (RecordForm, error) {
	var form RecordForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return RecordForm{}, err
	}
	form.ID = strings.TrimSpace(form.ID)
	form.ClientID = strings.TrimSpace(form.ClientID)
	form.ClientName = domain.NormalizeName(form.ClientName)
	form.ServiceTypeName = strings.TrimSpace(form.ServiceTypeName)
	form.CaregiverName = domain.NormalizeName(form.CaregiverName)
	return form, nil
}

func (f RecordForm) Validate() error {
	if f.ID == "" || f.ClientID == "" || f.ClientName == "" {
		return errors.New("id, client id and client name are required")
	}
	if f.VisitDate == "" || f.NextFollowUp == "" {
		return errors.New("visit dates are required")
	}
	if f.Score < 1 || f.Score > 5 {
		return errors.New("score must be between 1 and 5")
	}
	if strings.TrimSpace(f.Comment) == "" {
		return errors.New("comment is required")
	}
	return nil
}

func (f RecordForm) ToRecord() (domain.FollowUpRecord, error) {
	if err := f.Validate(); err != nil {
		return domain.FollowUpRecord{}, err
	}
	visit, err := parseDate(f.VisitDate)
	if err != nil {
		return domain.FollowUpRecord{}, err
	}
	next, err := parseDate(f.NextFollowUp)
	if err != nil {
		return domain.FollowUpRecord{}, err
	}
	return domain.FollowUpRecord{ID: f.ID, ClientID: f.ClientID, ClientName: f.ClientName, ServiceTypeID: f.ServiceTypeID, ServiceTypeName: f.ServiceTypeName, CaregiverID: f.CaregiverID, CaregiverName: f.CaregiverName, VisitDate: visit, NextFollowUp: next, Satisfaction: domain.Satisfaction{Score: f.Score, Comment: f.Comment}, Status: "scheduled", Notes: f.Notes, UpdatedAt: visit}, nil
}

func parseDate(value string) (time.Time, error) {
	return time.Parse(time.DateOnly, strings.TrimSpace(value))
}

func statusForError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if strings.Contains(err.Error(), "not found") {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

func allowMethods(w http.ResponseWriter, allowed ...string) {
	w.Header().Set("Allow", strings.Join(allowed, ", "))
}
