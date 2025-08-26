package database

import "time"

type AskForVip struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	UUID       string    `json:"uuid" `
	Message    string    `json:"message"`
	CreateTime time.Time `json:"create_time"`
	FinishTime time.Time `json:"finish_time"`
	Approver   string `json:"approver"`
	ApprovalResults string `json:"approval_results"`
}
