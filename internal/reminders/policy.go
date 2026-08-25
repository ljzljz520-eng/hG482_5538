package reminders

import (
	"fmt"
	"strings"
	"time"

	"homemaker-followup/internal/domain"
)

func NormalizeSetting(setting domain.ReminderSetting) (domain.ReminderSetting, error) {
	if setting.Channel == "" {
		setting.Channel = domain.ChannelPhone
	}
	if setting.QuietStart == setting.QuietEnd {
		setting.QuietStart = 21
		setting.QuietEnd = 8
	}
	if err := setting.Validate(); err != nil {
		return domain.ReminderSetting{}, err
	}
	return setting, nil
}

func CanSend(setting domain.ReminderSetting, at time.Time) bool {
	if !setting.Enabled {
		return false
	}
	hour := at.Hour()
	if setting.QuietStart < setting.QuietEnd {
		return hour < setting.QuietStart || hour >= setting.QuietEnd
	}
	return hour >= setting.QuietEnd && hour < setting.QuietStart
}

func RouteChannel(client domain.Client, setting domain.ReminderSetting) domain.Channel {
	if client.PreferredChannel != "" {
		return client.PreferredChannel
	}
	return setting.Channel
}

func ComposeMessage(reminder Reminder) string {
	if reminder.DaysUntil < 0 {
		return fmt.Sprintf("%s您好，您的家政服务回访已逾期，请联系客服安排时间。", reminder.ClientName)
	}
	if reminder.DaysUntil == 0 {
		return fmt.Sprintf("%s您好，今天是您的家政回访日。", reminder.ClientName)
	}
	return fmt.Sprintf("%s您好，距离家政回访还有%d天。", reminder.ClientName, reminder.DaysUntil)
}

func Snooze(reminder Reminder, days int) Reminder {
	if days < 0 {
		days = 0
	}
	reminder.DueDate = reminder.DueDate.AddDate(0, 0, days)
	reminder.DaysUntil += days
	return reminder
}

func Deduplicate(items []Reminder) []Reminder {
	seen := make(map[string]bool, len(items))
	result := make([]Reminder, 0, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item.RecordID)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, item)
	}
	return result
}
