package entity

import (
	"time"
)

const (
	MailTrackerStatusNew     = "new"
	MailTrackerStatusQueued  = "queued"
	MailTrackerStatusSuccess = "success"
	MailTrackerStatusError   = "error"
)

type mailTrackerStatus struct {
	MailTrackerStatusSuccess string
	MailTrackerStatusError   string
	MailTrackerStatusQueued  string
}

var MailTrackerStatusAll = mailTrackerStatus{
	MailTrackerStatusSuccess: MailTrackerStatusSuccess,
	MailTrackerStatusError:   MailTrackerStatusError,
	MailTrackerStatusQueued:  MailTrackerStatusQueued,
}

type MailTrackerEntity struct {
	ID           uint64 `orm:"table=email_tracker"`
	Status       string `orm:"enum=entity.MailTrackerStatusAll"`
	From         string `orm:"varchar=255"`
	To           string `orm:"varchar=255"`
	Subject      string
	TemplateFile string
	TemplateData string `orm:"length=max"`
	SenderError  string
	ReadAt       *time.Time `orm:"time"`
	CreatedAt    time.Time  `orm:"time"`
}
