package reporting

import (
	"fmt"
	"sort"
	"strings"

	"homemaker-followup/internal/domain"
)

type ServiceSummary struct {
	Service        string
	Records        int
	AverageScore   float64
	NeedsAttention int
}

func BuildServiceSummaries(records []domain.FollowUpRecord) []ServiceSummary {
	groups := make(map[string][]domain.FollowUpRecord)
	for _, record := range records {
		groups[record.ServiceTypeName] = append(groups[record.ServiceTypeName], record)
	}
	result := make([]ServiceSummary, 0, len(groups))
	for name, items := range groups {
		attention := 0
		for _, item := range items {
			if domain.SatisfactionBand(item.Satisfaction.Score).Label == "needs-attention" {
				attention++
			}
		}
		result = append(result, ServiceSummary{Service: name, Records: len(items), AverageScore: domain.AverageScore(items), NeedsAttention: attention})
	}
	sort.Slice(result, func(left, right int) bool { return result[left].Records > result[right].Records })
	return result
}

func RenderPlainText(summaries []ServiceSummary) string {
	lines := make([]string, 0, len(summaries)+1)
	lines = append(lines, "服务类型,记录数,平均满意度,待改进")
	for _, summary := range summaries {
		lines = append(lines, fmt.Sprintf("%s,%d,%.2f,%d", summary.Service, summary.Records, summary.AverageScore, summary.NeedsAttention))
	}
	return strings.Join(lines, "\n")
}

func SatisfactionDistribution(records []domain.FollowUpRecord) map[int]int {
	result := make(map[int]int)
	for _, record := range records {
		result[record.Satisfaction.Score]++
	}
	return result
}
