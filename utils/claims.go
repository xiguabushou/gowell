package utils

import (
	"goMedia/global"
	"goMedia/model/appTypes"
	"goMedia/model/request"
	"net"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAccessToken 从请求头获取Access Token
func GetAccessToken(c *gin.Context) string {
	// 获取x-access-token头部值
	token := c.Request.Header.Get("x-access-token")
	if len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"' {
		return token[1 : len(token)-1]
	}
	return token
}

// GetClaims 从Gin的Context中解析并获取JWT的Claims
func GetClaims(c *gin.Context) (*request.JwtCustomClaims, error) {
	// 获取请求头中的Access Token
	token := GetAccessToken(c)
	// 创建JWT实例
	j := NewJWT()
	// 解析Access Token
	claims, err := j.ParseAccessToken(token)
	if err != nil {
		// 如果解析失败，记录错误日志
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context. Please check if the request header contains 'x-access-token' and if the claims structure is correct.", zap.Error(err))
	}
	return claims, err
}

// GetUUID 从Gin的Context中获取JWT解析出来的用户UUID
func GetUUID(c *gin.Context) string {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回一个空UUID
			return ""
		} else {
			// 返回解析出来的UUID
			return cl.UUID
		}
	} else {
		// 如果已存在claims，则直接返回UUID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.UUID
	}
}

// GetRoleID 从Gin的Context中获取JWT解析出来的用户角色ID
func GetRoleID(c *gin.Context) appTypes.RoleID {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回0
			return 0
		} else {
			// 返回解析出来的角色ID
			return cl.RoleID
		}
	} else {
		// 如果已存在claims，则直接返回角色ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.RoleID
	}
}

// GetEmail 从Gin的Context中获取JWT解析出来的用户邮箱
func GetEmail(c *gin.Context) string {
	// 首先尝试从Context中获取"claims"
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析Access Token
		if cl, err := GetClaims(c); err != nil {
			// 如果解析失败，返回0
			return ""
		} else {
			// 返回解析出来的角色ID
			return cl.Email
		}
	} else {
		// 如果已存在claims，则直接返回角色ID
		waitUse := claims.(*request.JwtCustomClaims)
		return waitUse.Email
	}
}

// ClearAccessToken 清除access Token的cookie
func ClearAccessToken(c *gin.Context) {
	// 获取请求的host，如果失败则取原始请求host
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	// 调用setCookie设置cookie值为空并过期，删除refresh-token
	setCookie(c, "x-access-token", "", -1, host)
}

// setCookie 设置指定名称和值的cookie
func setCookie(c *gin.Context, name, value string, maxAge int, host string) {
	// 判断host是否是IP地址
	if net.ParseIP(host) != nil {
		// 如果是IP地址，设置cookie的domain为“/”
		c.SetCookie(name, value, maxAge, "/", "", false, true)
	} else {
		// 如果是域名，设置cookie的domain为域名
		c.SetCookie(name, value, maxAge, "/", host, false, true)
	}
}
