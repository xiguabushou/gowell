package database

import (
	"goMedia/model/appTypes"
	"time"
)

type Content struct{
	ID        uint `json:"id" gorm:"primarykey"`
	UID string `json:"uid"`
	TypeID appTypes.TypeID `json:"type_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title string `json:"title"`
	Tags string `json:"tags"`
	Number int `json:"number"`
}