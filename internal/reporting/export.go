package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"homemaker-followup/internal/domain"
)

type ExportBundle struct {
	Summaries    []ServiceSummary `json:"summaries"`
	Distribution map[int]int      `json:"distribution"`
	Records      int              `json:"records"`
}

func ExportJSON(path string, records []domain.FollowUpRecord) error {
	bundle := ExportBundle{Summaries: BuildServiceSummaries(records), Distribution: SatisfactionDistribution(records), Records: len(records)}
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func ValidateExport(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("report is empty")
	}
	return nil
}
