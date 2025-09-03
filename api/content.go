package api

import (
	"fmt"
	"goMedia/global"
	"goMedia/model/request"
	"goMedia/model/response"

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
	res, err := contentService.GetInfo(req.UID)
	if err != nil{
		global.Log.Error("Failed to get content info:", zap.Error(err))
		response.FailWithMessage("Failed to get content info", c)
		return
	}

	response.OkWithData(res,c)
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
	err := c.ShouldBindQuery(&req)
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

func (contentApi *ContentApi)Freeze(c *gin.Context){
	var req request.GetID
	err:= c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.Freeze(req.UID);err != nil{
		global.Log.Error("Failed to get freeze:", zap.Error(err))
		response.FailWithMessage("Failed to get freeze", c)
		return
	}
	response.Ok(c)
}

func (contentApi *ContentApi)UnFreeze(c *gin.Context){
	var req request.GetID
	err:= c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.UnFreeze(req.UID);err != nil{
		global.Log.Error("Failed to unfreeze:", zap.Error(err))
		response.FailWithMessage("Failed to unfreeze", c)
		return
	}
	response.Ok(c)
}

func (contentApi *ContentApi)Delete(c *gin.Context){
	var req request.GetID
	err:= c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := contentService.Delete(req.UID);err != nil{
		global.Log.Error("Failed to delete:", zap.Error(err))
		response.FailWithMessage("Failed to delete", c)
		return
	}
	response.Ok(c)
}