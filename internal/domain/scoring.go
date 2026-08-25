package domain

import "sort"

type ScoreBand struct {
	Label    string
	Minimum  int
	Maximum  int
	Priority int
}

var scoreBands = []ScoreBand{
	{Label: "needs-attention", Minimum: 1, Maximum: 2, Priority: 3},
	{Label: "stable", Minimum: 3, Maximum: 4, Priority: 2},
	{Label: "advocate", Minimum: 5, Maximum: 5, Priority: 1},
}

func SatisfactionBand(score int) ScoreBand {
	for _, band := range scoreBands {
		if score >= band.Minimum && score <= band.Maximum {
			return band
		}
	}
	return ScoreBand{Label: "unknown", Priority: 4}
}

func PriorityForRecord(record FollowUpRecord) int {
	priority := SatisfactionBand(record.Satisfaction.Score).Priority
	if record.NextFollowUp.IsZero() {
		return priority + 2
	}
	if record.Status == "overdue" {
		return priority + 3
	}
	return priority
}

func RankRecords(records []FollowUpRecord) []FollowUpRecord {
	result := append([]FollowUpRecord(nil), records...)
	sort.SliceStable(result, func(left, right int) bool {
		leftPriority := PriorityForRecord(result[left])
		rightPriority := PriorityForRecord(result[right])
		if leftPriority != rightPriority {
			return leftPriority > rightPriority
		}
		return result[left].NextFollowUp.Before(result[right].NextFollowUp)
	})
	return result
}

func AverageScore(records []FollowUpRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	total := 0
	for _, record := range records {
		total += record.Satisfaction.Score
	}
	return float64(total) / float64(len(records))
}

func CountByStatus(records []FollowUpRecord) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		counts[DisplayStatus(record)]++
	}
	return counts
}
