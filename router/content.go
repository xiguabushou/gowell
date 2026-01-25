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
		contentVipRouter.GET("getList", contentApi.GetList) //用户获取内容列表
		contentVipRouter.GET("getInfo", contentApi.GetInfo) //用户获取内容详情
	}
	{
		contentAdminRouter.POST("uploadVideo", contentApi.UploadVideo)                   //上传视频内容
		contentAdminRouter.POST("uploadPhoto", contentApi.UploadPhoto)                   //上传图片内容
		contentAdminRouter.POST("listByAdmin", contentApi.ListByAdmin)                   //获取所有内容列表 (包括已下架内容)
		contentAdminRouter.POST("freeze", contentApi.Freeze)                             //下架/取消下架内容
		contentAdminRouter.POST("delete", contentApi.Delete)                             //删除内容(整体删除)
		contentAdminRouter.POST("editContentPhotoInfo", contentApi.EditContentPhotoInfo) //编辑标题和标签
		contentAdminRouter.POST("uploadContentVideo", contentApi.UploadContentVideo)     //上传视频到原有内容中
		contentAdminRouter.POST("deleteContentPhoto", contentApi.DeleteContentPhoto)     //删除原有内容中的图片
		contentAdminRouter.POST("uploadContentPhoto", contentApi.UploadContentPhoto)     //上传图片到原有的内容中
		contentAdminRouter.POST("editVideo", contentApi.EditVideo)                       //获取要编辑的视频相关信息
		contentAdminRouter.POST("editPhoto", contentApi.EditPhoto)                       //获取要编辑图片的相关信息
	}
}
