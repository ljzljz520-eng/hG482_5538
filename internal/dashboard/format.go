package dashboard

import (
	"fmt"
	"sort"
	"strings"

	"homemaker-followup/internal/domain"
)

type TableRow struct {
	ID        string `json:"id"`
	Client    string `json:"client"`
	Service   string `json:"service"`
	Caregiver string `json:"caregiver"`
	Status    string `json:"status"`
	Score     int    `json:"score"`
	Next      string `json:"next"`
}

func ToTableRows(records []domain.FollowUpRecord) []TableRow {
	result := make([]TableRow, 0, len(records))
	for _, record := range records {
		result = append(result, TableRow{ID: record.ID, Client: record.ClientName, Service: record.ServiceTypeName, Caregiver: record.CaregiverName, Status: domain.DisplayStatus(record), Score: record.Satisfaction.Score, Next: domain.FormatDate(record.NextFollowUp)})
	}
	sort.SliceStable(result, func(left, right int) bool { return result[left].Next < result[right].Next })
	return result
}

func MetricLabels(snapshot Snapshot) map[string]string {
	return map[string]string{"clients": fmt.Sprintf("%d", snapshot.TotalClients), "records": fmt.Sprintf("%d", snapshot.TotalRecords), "average": fmt.Sprintf("%.2f", snapshot.AverageScore), "attention": fmt.Sprintf("%d", snapshot.StatusCounts["needs-attention"])}
}

func CalendarLabels(snapshot Snapshot) []string {
	labels := make([]string, 0, len(snapshot.Calendar))
	for _, day := range snapshot.Calendar {
		labels = append(labels, fmt.Sprintf("%s:%d", day.Date.Format("01-02"), day.Total))
	}
	return labels
}

func SearchHints(records []domain.FollowUpRecord) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, record := range records {
		for _, value := range []string{record.ClientName, record.CaregiverName, record.ServiceTypeName} {
			value = strings.TrimSpace(value)
			if value != "" && !seen[value] {
				seen[value] = true
				result = append(result, value)
			}
		}
	}
	sort.Strings(result)
	return result
}

func HealthLabel(snapshot Snapshot) string {
	if snapshot.TotalRecords == 0 {
		return "等待第一条回访"
	}
	if snapshot.StatusCounts["overdue"] > 0 {
		return "存在逾期回访"
	}
	if snapshot.StatusCounts["needs-attention"] > 0 {
		return "存在待改进客户"
	}
	return "回访状态稳定"
}
