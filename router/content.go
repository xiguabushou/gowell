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
		contentVipRouter.GET("home", contentApi.Home)
		contentVipRouter.GET("video", contentApi.Video)
		contentVipRouter.GET("photo", contentApi.Photo)
		contentVipRouter.GET("search", contentApi.Search)
	}
	{
		contentAdminRouter.POST("uploadVideo", contentApi.UploadVideo)
		contentAdminRouter.POST("uploadPhoto", contentApi.UploadPhoto)
		contentAdminRouter.POST("editVideo", contentApi.EditVideo)
		contentAdminRouter.POST("editPhoto", contentApi.EditPhoto)
		contentAdminRouter.POST("list", contentApi.List)
	}
}
