package middleware

import (
	"fmt"
	"goMedia/model/appTypes"
	"goMedia/model/response"
	"goMedia/utils"

	"github.com/gin-gonic/gin"
)

func VipAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleId := utils.GetRoleID(c)
		fmt.Println(roleId)

		if roleId != appTypes.Admin && roleId != appTypes.Vip {
			response.Forbidden("Access denied. vip privileges are required", c)
			c.Abort()
			return
		}

		c.Next()
	}
}
