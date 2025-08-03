package database

import "goMedia/global"

// Login 登录日志表
type Login struct {
	global.MODEL
	UserID      string `json:"user_id"`      // 用户 ID
	Email       string `json:"email"`        //邮箱
	IP          string `json:"ip"`           // IP 地址
	Address     string `json:"address"`      // 登录地址
	OS          string `json:"os"`           // 操作系统
	DeviceInfo  string `json:"device_info"`  // 设备信息
	BrowserInfo string `json:"browser_info"` // 浏览器信息
	Status      int    `json:"status"`       // 登录状态
}
