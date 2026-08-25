package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/excel"
	"homemaker-followup/internal/followup"
	"homemaker-followup/internal/httpapi"
	"homemaker-followup/internal/storage"
)

func main() {
	dataDir := flag.String("data", "./data", "数据目录")
	addr := flag.String("addr", ":8080", "HTTP监听地址")
	flag.Parse()
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatal(err)
	}
	workbookPath := filepath.Join(*dataDir, "回访簿.xlsx")
	if _, err := os.Stat(workbookPath); os.IsNotExist(err) {
		workbook, createErr := excel.NewWorkbook(workbookPath)
		if createErr != nil {
			log.Fatal(createErr)
		}
		if saveErr := workbook.Save(); saveErr != nil {
			log.Fatal(saveErr)
		}
		workbook.Close()
	}
	store, err := storage.Open(filepath.Join(*dataDir, "followup.sqlite"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	service := followup.NewService(workbookPath, store)
	notice, startErr := service.Start()
	fmt.Printf("%s: %s\n", notice.Title, notice.Detail)
	if startErr != nil && notice.Blocking {
		log.Fatal(startErr)
	}
	server := httpapi.NewServer(service, seedClients(), seedSetting(), time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
	log.Printf("家政客户回访簿运行于 %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}

func seedClients() []domain.Client {
	return []domain.Client{{ID: "C-001", Name: "林女士", Phone: "13800000001", Address: "静安区", PreferredChannel: domain.ChannelWeChat, Active: true}}
}

func seedSetting() domain.ReminderSetting {
	return domain.ReminderSetting{ID: "default", DaysBefore: 3, Channel: domain.ChannelWeChat, Enabled: true, QuietStart: 21, QuietEnd: 8}
}
