package followup

import (
	"context"
	"sort"

	"homemaker-followup/internal/domain"
)

func (s *Service) Search(query string) []domain.FollowUpRecord {
	query = normalizeQuery(query)
	result := make([]domain.FollowUpRecord, 0)
	for _, record := range s.Records {
		if query == "" || containsRecord(record, query) {
			result = append(result, record)
		}
	}
	sort.SliceStable(result, func(left, right int) bool { return result[left].NextFollowUp.Before(result[right].NextFollowUp) })
	return result
}

func (s *Service) RefreshFromIndex(ctx context.Context) error {
	if s.Store == nil {
		return nil
	}
	records, err := s.Store.ListRecords(ctx)
	if err != nil {
		return err
	}
	s.Records = records
	return nil
}

func normalizeQuery(value string) string {
	return domain.NormalizeName(value)
}

func containsRecord(record domain.FollowUpRecord, query string) bool {
	fields := []string{record.ID, record.ClientID, record.ClientName, record.ServiceTypeName, record.CaregiverName, record.Status, record.Notes}
	for _, field := range fields {
		if normalizeQuery(field) == query || containsFold(field, query) {
			return true
		}
	}
	return false
}

func containsFold(value, query string) bool {
	valueLower := normalizeQuery(value)
	queryLower := normalizeQuery(query)
	for index := 0; index+len(queryLower) <= len(valueLower); index++ {
		if valueLower[index:index+len(queryLower)] == queryLower {
			return true
		}
	}
	return false
}
