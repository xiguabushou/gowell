package database

import "time"

type AskForVip struct {
	ID              uint   `json:"id" gorm:"primarykey"`
	UUID            string `json:"uuid" `
	Message         string `json:"message"`
	CreatedAt       time.Time
	FinishAt       time.Time `json:"finish_at"`
	Approver        string `json:"approver"`
	ApprovalResults string `json:"approval_results"`
}
