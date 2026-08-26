package reminders

import (
	"fmt"
	"time"

	"homemaker-followup/internal/domain"
)

type CalendarDay struct {
	Date           time.Time
	Total          int
	Overdue        int
	NeedsAttention int
}

func SevenDayCalendar(records []domain.FollowUpRecord, start time.Time) []CalendarDay {
	result := make([]CalendarDay, 7)
	for index := range result {
		result[index].Date = start.AddDate(0, 0, index)
	}
	for _, record := range records {
		for index := range result {
			if sameDate(record.NextFollowUp, result[index].Date) {
				result[index].Total++
				if record.Status == "overdue" {
					result[index].Overdue++
				}
				if domain.SatisfactionBand(record.Satisfaction.Score).Label == "needs-attention" {
					result[index].NeedsAttention++
				}
			}
		}
	}
	return result
}

func sameDate(left, right time.Time) bool {
	return left.Year() == right.Year() && left.YearDay() == right.YearDay()
}

func DescribeDay(day CalendarDay) string {
	if day.Total == 0 {
		return fmt.Sprintf("%s: 无回访", day.Date.Format(time.DateOnly))
	}
	if day.Overdue > 0 {
		return fmt.Sprintf("%s: %d项，其中%d项逾期", day.Date.Format(time.DateOnly), day.Total, day.Overdue)
	}
	return fmt.Sprintf("%s: %d项待回访", day.Date.Format(time.DateOnly), day.Total)
}
