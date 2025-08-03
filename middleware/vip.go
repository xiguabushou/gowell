package middleware

import (
	"github.com/gin-gonic/gin"
	"goMedia/model/appTypes"
	"goMedia/model/response"
	"goMedia/utils"
)

func VipAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleId := utils.GetRoleID(c)

		if roleId == appTypes.Admin || roleId == appTypes.Vip {
			response.Forbidden("Access denied. vip privileges are required", c)
			c.Abort()
			return
		}

		c.Next()
	}
}
