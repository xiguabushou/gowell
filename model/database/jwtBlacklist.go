package database

import (
	"time"
)

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	ID          uint      `json:"id" gorm:"primarykey"` // 主键 ID
	CreatedTime time.Time `json:"created_Time"`         // 创建时间
	Jwt         string    `json:"jwt" gorm:"type:text"` // Jwt
}
