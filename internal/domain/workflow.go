package domain

import (
	"errors"
	"fmt"
	"time"
)

type Intake struct {
	Client       Client
	Service      ServiceType
	Caregiver    Caregiver
	VisitDate    time.Time
	NextFollowUp time.Time
	Score        int
	Comment      string
	Notes        string
}

func (i Intake) BuildRecord(id string) (FollowUpRecord, error) {
	if err := i.Client.Validate(); err != nil {
		return FollowUpRecord{}, fmt.Errorf("client: %w", err)
	}
	if err := i.Service.Validate(); err != nil {
		return FollowUpRecord{}, fmt.Errorf("service: %w", err)
	}
	if err := i.Caregiver.Validate(); err != nil {
		return FollowUpRecord{}, fmt.Errorf("caregiver: %w", err)
	}
	if i.VisitDate.IsZero() {
		return FollowUpRecord{}, errors.New("visit date is required")
	}
	next := i.NextFollowUp
	if next.IsZero() {
		next = i.VisitDate.AddDate(0, 0, i.Service.DefaultDays)
	}
	record := FollowUpRecord{
		ID: id, ClientID: i.Client.ID, ClientName: NormalizeName(i.Client.Name),
		ServiceTypeID: i.Service.ID, ServiceTypeName: i.Service.Name,
		CaregiverID: i.Caregiver.ID, CaregiverName: NormalizeName(i.Caregiver.Name),
		VisitDate: i.VisitDate, NextFollowUp: next,
		Satisfaction: Satisfaction{Score: i.Score, Comment: i.Comment}, Notes: i.Notes,
		Status: "scheduled", UpdatedAt: i.VisitDate.UTC(),
	}
	if err := record.Validate(); err != nil {
		return FollowUpRecord{}, err
	}
	return record, nil
}

func MarkOverdue(record FollowUpRecord, today time.Time) FollowUpRecord {
	if !record.NextFollowUp.IsZero() && record.NextFollowUp.Before(today) && record.Status == "scheduled" {
		record.Status = "overdue"
	}
	return record
}

func CompleteFollowUp(record FollowUpRecord, score int, comment string, next time.Time) (FollowUpRecord, error) {
	if score < 1 || score > 5 {
		return FollowUpRecord{}, errors.New("score must be between 1 and 5")
	}
	if next.Before(record.VisitDate) {
		return FollowUpRecord{}, errors.New("next follow-up cannot precede visit")
	}
	record.Satisfaction = Satisfaction{Score: score, Comment: comment}
	record.NextFollowUp = next
	record.Status = "scheduled"
	record.UpdatedAt = next.UTC()
	return record, record.Validate()
}
