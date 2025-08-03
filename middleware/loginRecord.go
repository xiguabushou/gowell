package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"go.uber.org/zap"
	"goMedia/global"
	"goMedia/model/database"
)

// LoginRecord 是一个中间件，用于记录登录日志
func LoginRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 异步记录日志
		go func() {
			var userID string
			var email string
			var address string
			ip := c.ClientIP()
			userAgent := c.Request.UserAgent()

			// 从请求上下文中获取用户ID，确保获取到的是当前请求的正确用户ID
			if value, exists := c.Get("user_id"); exists {
				if id, ok := value.(string); ok {
					userID = id
				}
			}

			if value, exists := c.Get("email"); exists {
				if id, ok := value.(string); ok {
					email = id
				}
			}

			// TODO 获取用户IP的地理位置

			// 解析用户的浏览器、操作系统和设备信息
			os, deviceInfo, browserInfo := parseUserAgent(userAgent)

			// 创建登录记录
			login := database.Login{
				UserID:      userID,
				Email:       email,
				IP:          ip,
				Address:     address,
				OS:          os,
				DeviceInfo:  deviceInfo,
				BrowserInfo: browserInfo,
				Status:      c.Writer.Status(),
			}

			// 将登录记录存储到数据库
			if err := global.DB.Create(&login).Error; err != nil {
				global.Log.Error("Failed to record login", zap.Error(err))
			}

		}()
	}
}

// 解析用户代理（User-Agent）字符串，提取操作系统、设备信息和浏览器信息
func parseUserAgent(userAgent string) (os, deviceInfo, browserInfo string) {
	os = userAgent
	deviceInfo = userAgent
	browserInfo = userAgent

	parser := uaparser.NewFromSaved()
	cli := parser.Parse(userAgent)
	os = cli.Os.Family
	deviceInfo = cli.Device.Family
	browserInfo = cli.UserAgent.Family

	return
}
