package database

import (
	"goMedia/model/appTypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     string          `json:"uuid" gorm:"type:char(36);unique"` //uuid
	Password string          `json:"-"`                                // 密码
	Email    string          `json:"email"`                            //邮箱
	RoleID   appTypes.RoleID `json:"role_id"`                          // 角色 ID
	Freeze   bool            `json:"freeze"`                           // 用户是否被冻结
}
