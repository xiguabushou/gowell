package router

import (
	"goMedia/api"
	"goMedia/middleware"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (b *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	pingRouter := Router.Group("base").Use(middleware.JWTAuth())

	baseApi := api.ApiGroupApp.BaseApi
	{
		pingRouter.GET("ping", baseApi.Ping)
		baseRouter.GET("captcha", baseApi.Captcha)
		baseRouter.POST("sendEmailVerificationCode", baseApi.SendEmailVerificationCode)
	}

}
