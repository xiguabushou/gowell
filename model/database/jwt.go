package database

import "time"

type Jwt struct {
	UUID        string    `json:"uuid" gorm:"primarykey"` //user uuid
	UpdatedTime time.Time `json:"updated_time"`           // 更新时间
	Jwt         string    `json:"jwt" gorm:"type:text"`   // Jwt
}
