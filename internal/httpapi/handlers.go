package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"homemaker-followup/internal/domain"
)

func decodeRecord(r *http.Request) (domain.FollowUpRecord, error) {
	var record domain.FollowUpRecord
	err := json.NewDecoder(r.Body).Decode(&record)
	return record, err
}

func parseID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func clientSummary(clients []domain.Client) map[string]domain.Client {
	result := make(map[string]domain.Client, len(clients))
	for _, client := range clients {
		result[client.ID] = client
	}
	return result
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
