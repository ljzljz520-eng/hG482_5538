package dashboard

import (
	"strings"

	"homemaker-followup/internal/domain"
)

type Filter struct {
	Status       string
	Caregiver    string
	ServiceType  string
	MinimumScore int
}

func ApplyFilter(records []domain.FollowUpRecord, filter Filter) []domain.FollowUpRecord {
	result := make([]domain.FollowUpRecord, 0)
	for _, record := range records {
		if filter.Status != "" && domain.DisplayStatus(record) != filter.Status {
			continue
		}
		if filter.Caregiver != "" && !strings.Contains(strings.ToLower(record.CaregiverName), strings.ToLower(filter.Caregiver)) {
			continue
		}
		if filter.ServiceType != "" && !strings.Contains(strings.ToLower(record.ServiceTypeName), strings.ToLower(filter.ServiceType)) {
			continue
		}
		if filter.MinimumScore > 0 && record.Satisfaction.Score < filter.MinimumScore {
			continue
		}
		result = append(result, record)
	}
	return result
}

func DistinctCaregivers(records []domain.FollowUpRecord) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, record := range records {
		if !seen[record.CaregiverName] {
			seen[record.CaregiverName] = true
			result = append(result, record.CaregiverName)
		}
	}
	return result
}
