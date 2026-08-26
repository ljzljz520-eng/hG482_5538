package followup

import (
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestFollowupHeaderError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.xlsx")
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "回访记录")
	file.SetCellValue("回访记录", "A1", "客户")
	if err := file.SaveAs(path); err != nil {
		t.Fatal(err)
	}
	file.Close()
	if _, err := LoadFromFile(path); err == nil {
		t.Fatal("expected startup format error")
	}
}
