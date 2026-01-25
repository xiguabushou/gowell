package other

import (
	"goMedia/model/appTypes"
	"time"
)

type ContentList struct {
	Uid   string   `json:"uid"`
	Title string   `json:"title"`
	Cover string   `json:"cover"`
	Tags  []string `json:"tags"`
}

type ListByAdmin struct {
	ID        uint            `json:"id" gorm:"primarykey"`
	UID       string          `json:"uid"`
	TypeID    appTypes.TypeID `json:"type_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string   `json:"title"`
	Tags      string `json:"tags" gorm:"type:json"`
	Freeze    bool     `json:"freeze"`
	Cover     string   `json:"cover"`
}
