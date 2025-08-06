package request

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"goMedia/model/appTypes"
)

// JwtCustomClaims 结构体用于存储JWT的自定义Claims，继承自BaseClaims，并包含标准的JWT注册信息
type JwtCustomClaims struct {
	BaseClaims           // 基础Claims，包含用户ID、UUID和角色ID
	jwt.RegisteredClaims // 标准JWT声明，例如过期时间、发行者等
}

type JwtCustomClaims2 struct {
	ForgotPasswordClaims // 基础Claims，包含用户ID、UUID和角色ID
	jwt.RegisteredClaims // 标准JWT声明，例如过期时间、发行者等
}

// BaseClaims 结构体用于存储基本的用户信息，作为JWT的Claim部分
type BaseClaims struct {
	UUID   string          // 用户ID，标识用户唯一性
	RoleID appTypes.RoleID // 用户角色ID，表示用户的权限级别
}

type ForgotPasswordClaims struct {
	Email string `json:"email"`
}
