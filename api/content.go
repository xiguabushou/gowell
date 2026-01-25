package api

import (
	"errors"
	"fmt"
	"goMedia/global"
	"goMedia/model/request"
	"goMedia/model/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ContentApi struct{}

func (contentApi *ContentApi) GetList(c *gin.Context) {
	var req request.GetList
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := contentService.GetList(req)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (contentApi *ContentApi) GetInfo(c *gin.Context) {
	var req request.GetInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	fmt.Println(req.UID)
	res, err := contentService.GetInfo(req.UID, req.Page, req.PageSize)
	if err != nil {
		global.Log.Error("Failed to get content info:", zap.Error(err))
		response.FailWithMessage("Failed to get content info", c)
		return
	}

	response.OkWithData(res, c)
}

func (contentApi *ContentApi) UploadVideo(c *gin.Context) {

	title := c.PostForm("title")
	tags := c.PostForm("tags")
	cover, err := c.FormFile("cover")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	file, err := c.FormFile("video")

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = contentService.UploadVideo(title, tags, file, cover, c)
	if err != nil {
		global.Log.Error("Failed to upload video:", zap.Error(err))
		response.FailWithMessage("Failed to upload video", c)
		return
	}
	response.OkWithMessage("Successfully to upload video", c)
}

func (contentApi *ContentApi) UploadPhoto(c *gin.Context) {
	title := c.PostForm("title")
	tags := c.PostForm("tags")
	cover, err := c.FormFile("cover")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	formdata := c.Request.MultipartForm
	files := formdata.File["photo"]

	err = contentService.UploadPhoto(title, tags, files, cover, c)
	if err != nil {
		global.Log.Error("Failed to upload photos:", zap.Error(err))
		response.FailWithMessage("Failed to upload photos", c)
		return
	}
	response.OkWithMessage("Successfully to upload photos", c)

}

func (contentApi *ContentApi) ListByAdmin(c *gin.Context) {
	var req request.ListByAdmin
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := contentService.ListByAdmin(req)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (contentApi *ContentApi) Freeze(c *gin.Context) {
	var req request.ContentFreeze
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.Freeze(req); err != nil {
		global.Log.Error("Failed to get freeze:", zap.Error(err))
		response.FailWithMessage("Failed to get freeze", c)
		return
	}
	response.Ok(c)
}

func (contentApi *ContentApi) Delete(c *gin.Context) {
	var req request.GetID
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.Delete(req.UID); err != nil {
		global.Log.Error("Failed to delete:", zap.Error(err))
		response.FailWithMessage("Failed to delete", c)
		return
	}
	response.Ok(c)
}

func (contentApi *ContentApi) EditContentPhotoInfo(c *gin.Context) {
	uid := c.PostForm("uid")
	title := c.PostForm("title")
	tags := c.PostForm("tags")
	cover, err := c.FormFile("cover")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.EditContentPhotoInfo(uid, title, tags, cover); err != nil {
		global.Log.Error("Failed to edit the info", zap.Error(err))
		response.FailWithMessage("修改信息失败", c)
		return
	}
	response.Ok(c)
}

func (contentApi *ContentApi) UploadContentVideo(c *gin.Context) {
	uid := c.PostForm("uid")
	title := c.PostForm("title")
	tags := c.PostForm("tags")
	cover, err := c.FormFile("cover")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		response.FailWithMessage(err.Error(), c)
		return
	}

	file, err := c.FormFile("video")

	if err != nil && !errors.Is(err, http.ErrMissingFile){
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = contentService.UploadContentVideo(uid, title, tags, file, cover, c)
	if err != nil {
		global.Log.Error("Failed to upload video:", zap.Error(err))
		response.FailWithMessage("Failed to upload video", c)
		return
	}
	response.OkWithMessage("Successfully to upload video", c)
}

func (contentApi *ContentApi) DeleteContentPhoto(c *gin.Context) {
	var req request.DeleteContentPhoto
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.DeleteContentPhoto(req); err != nil {
		global.Log.Error("Failed to delete photo", zap.Error(err))
		response.FailWithMessage("删除图片失败", c)
		return
	}
	response.Ok(c)

}

func (contentApi *ContentApi) UploadContentPhoto(c *gin.Context) {
	uid := c.PostForm("uid")
	formdata := c.Request.MultipartForm
	files := formdata.File["photo"]

	err := contentService.UploadContentPhoto(uid,  files, c)
	if err != nil {
		global.Log.Error("Failed to upload photos:", zap.Error(err))
		response.FailWithMessage("Failed to upload photos", c)
		return
	}
	response.OkWithMessage("Successfully to upload photos", c)
}

func (contentApi *ContentApi) EditVideo(c *gin.Context) {
	var req request.GetID
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	videoInfo, err := contentService.EditVideo(req.UID)
	if err != nil {
		global.Log.Error("Failed to get edit video info", zap.Error(err))
		response.FailWithMessage("获取编辑视频信息失败", c)
		return
	}
	response.OkWithData(videoInfo, c)
}

func (contentApi *ContentApi) EditPhoto(c *gin.Context) {
	var req request.GetInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	photoInfo, err := contentService.EditPhoto(req)
	if err != nil {
		global.Log.Error("Failed to get edit photo info", zap.Error(err))
		response.FailWithMessage("获取编辑图片信息失败", c)
		return
	}
	response.OkWithData(photoInfo, c)
}
