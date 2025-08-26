package router

import (
	"github.com/gin-gonic/gin"
	"goMedia/api"
	"goMedia/middleware"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(PublicRouter *gin.RouterGroup, UserRouter *gin.RouterGroup, VipRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())
	userPublicRouter := PublicRouter.Group("user")
	userUserRouter := UserRouter.Group("user")
	userAdminRouter := AdminRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi
	{
		userLoginRouter.POST("register", userApi.Register)
		userLoginRouter.POST("login", userApi.Login)
	}
	{
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword)
		userPublicRouter.POST("resetForgotPassword", userApi.ResetForgotPassword)
	}
	{
		userUserRouter.POST("askForVip", userApi.AskForVip)
		userUserRouter.POST("logout", userApi.Logout)
		userUserRouter.POST("resetPassword", userApi.UserResetPassword)
		userUserRouter.GET("info", userApi.UserInfo)
	}
	{
		userAdminRouter.POST("edit", userApi.EditUser)
		userAdminRouter.POST("add", userApi.AddUser)
		userAdminRouter.POST("delete", userApi.DeleteUser)
		userAdminRouter.GET("list", userApi.UserList)
		userAdminRouter.POST("freeze", userApi.UserFreeze)
		userAdminRouter.POST("unfreeze", userApi.UserUnfreeze)
		userAdminRouter.GET("loginList", userApi.UserLoginList)
	}

}
