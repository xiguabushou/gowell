package router

import (
	"goMedia/api"

	"github.com/gin-gonic/gin"
)

type ContentRouter struct{}

func (c ContentRouter) InitContentRouter(VipRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	contentVipRouter := VipRouter.Group("content")
	contentAdminRouter := VipRouter.Group("content")
	contentApi := api.ApiGroupApp.ContentApi
	{
		contentVipRouter.GET("getList", contentApi.GetList)
		contentVipRouter.GET("getInfo", contentApi.GetInfo)
	}
	{
		contentAdminRouter.POST("uploadVideo", contentApi.UploadVideo)
		contentAdminRouter.POST("uploadPhoto", contentApi.UploadPhoto)
		contentAdminRouter.POST("listByAdmin", contentApi.ListByAdmin)
	}
}
