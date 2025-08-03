package middleware

import (
	"github.com/gin-gonic/gin"
	"goMedia/model/appTypes"
	"goMedia/model/response"
	"goMedia/utils"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleId := utils.GetRoleID(c)

		if roleId == appTypes.Admin {
			c.Next()
		}

		response.Forbidden("Access denied. Admin privileges are required", c)
		c.Abort()
		return
	}
}
