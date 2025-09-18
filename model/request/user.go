package request

import (
	"goMedia/model/appTypes"
	"time"
)

type Register struct {
	Password         string `json:"password" binding:"required,min=6,max=16"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
}

type Login struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=16"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

type ForgotPassword struct {
	Email     string `json:"email" binding:"required,email"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

type UserResetPassword struct {
	UUID        string `json:"-"`
	Password    string `json:"password" binding:"required,min=8,max=16"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=16"`
}

type UserOperation struct {
	ID string `json:"id" binding:"required"`
}

type EditUser struct {
	UUID     string          `json:"uuid" binding:"required"`
	Password string          `json:"password"`
	RoleID   appTypes.RoleID `json:"role_id"`
	Freeze   bool            `json:"freeze"`
	Email    string          `json:"email"`
}

type AddUser struct {
	Password string          `json:"password" binding:"required,min=8,max=16"`
	RoleID   appTypes.RoleID `json:"role_id"`
	Freeze   bool            `json:"freeze"`
	Email    string          `json:"email" binding:"required ,email"`
}

type UserList struct {
	Search string `json:"search"`
	PageInfo
}

type UserLoginList struct {
	Search string `json:"search"`
	PageInfo
}

type ResetForgotPassword struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type AskForVip struct {
	Message string `json:"message" binding:"max=150"`
	UUID    string `json:"uuid" binding:"required"`
}

type GetListAboutAskForVip struct {
	UUID      string    `json:"uuid" `
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type ApprovingForVip struct {
	ApproverUUID string `json:"approver_uuid"`
	UUID    string `json:"uuid"`
	IsPass  bool   `json:"is_pass"`
}
