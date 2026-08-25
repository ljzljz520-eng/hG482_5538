package domain

import (
	"fmt"
	"strings"
	"unicode"
)

func NormalizeName(value string) string {
	words := strings.Fields(strings.TrimSpace(value))
	for index, word := range words {
		runes := []rune(strings.ToLower(word))
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[index] = string(runes)
	}
	return strings.Join(words, " ")
}

func NormalizePhone(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if unicode.IsDigit(r) || r == '+' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func NormalizeTags(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, value := range values {
		tag := strings.ToLower(strings.TrimSpace(value))
		if tag != "" && !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}
	return result
}

func DisplayStatus(record FollowUpRecord) string {
	if record.Satisfaction.Score <= 2 {
		return "needs-attention"
	}
	if record.Status != "" {
		return record.Status
	}
	if record.NextFollowUp.IsZero() {
		return "unscheduled"
	}
	return "scheduled"
}

func FormatDate(value interface{ Format(string) string }) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}

func RecordKey(clientID, visitDate string) string {
	return fmt.Sprintf("%s:%s", strings.TrimSpace(clientID), strings.TrimSpace(visitDate))
}
