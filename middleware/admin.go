package middleware

import (
	"goMedia/model/appTypes"
	"goMedia/model/response"
	"goMedia/utils"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleId := utils.GetRoleID(c)

		if roleId == appTypes.Admin {
			c.Next()
			return
		}

		response.Forbidden("Access denied. Admin privileges are required", c)
		c.Abort()
	}
}
