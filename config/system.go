package config

import (
	"fmt"
)

// System 系统配置
type System struct {
	Host              string `json:"-" yaml:"host"`                                  // 服务器绑定的主机地址，通常为 0.0.0.0 表示监听所有可用地址
	Port              int    `json:"-" yaml:"port"`                                  // 服务器监听的端口号，通常用于 HTTP 服务
	Env               string `json:"-" yaml:"env"`                                   // Gin 的环境类型，例如 "debug"、"release" 或 "test"
	RouterPrefix      string `json:"-" yaml:"router_prefix"`                         // API 路由前缀，用于构建 API 路径
	UseMultipoint     bool   `json:"use_multipoint" yaml:"use_multipoint"`           // 是否启用多点登录拦截，防止同一账户在多个地方同时登录
	SessionsSecret    string `json:"sessions_secret" yaml:"sessions_secret"`         // 用于加密会话的密钥，确保会话数据的安全性
	ForgotPasswordUrl string `json:"forgot_password_url" yaml:"forgot_password_url"` //用于重新设置密码
}

func (s System) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
