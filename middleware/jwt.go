package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goMedia/model/response"
	"goMedia/service"
	"goMedia/utils"
)

var jwtService = service.ServiceGroupApp.JwtService

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := utils.GetAccessToken(c)

		if jwtService.IsInBlacklist(accessToken) {
			utils.ClearAccessToken(c)
			response.NoAuth("Account logged in from another location or token is invalid", c)
			c.Abort()
			return
		}

		j := utils.NewJWT()

		claims, err := j.ParseAccessToken(accessToken)
		if err != nil {
			if accessToken == "" || errors.Is(err, utils.TokenExpired) {
				utils.ClearAccessToken(c)
				response.NoAuth("access token expired or invalid", c)
				c.Abort()
				return
			}
			utils.ClearAccessToken(c)
			response.NoAuth("Parse access token error", c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
