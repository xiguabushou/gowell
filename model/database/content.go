package database

import (
	"goMedia/model/appTypes"
	"time"
)

type Content struct {
	ID        uint            `json:"id" gorm:"primarykey"`
	UID       string          `json:"uid"`
	TypeID    appTypes.TypeID `json:"type_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `json:"title"`
	Tags      []byte `json:"tags" gorm:"type:json"`
	Number    int    `json:"number"`
	Freeze    bool   `json:"freeze"`
}
