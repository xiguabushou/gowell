package api

import (
	"goMedia/global"
	"goMedia/model/request"
	"goMedia/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ContentApi struct{}

func (contentApi *ContentApi) Home(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := contentService.Home(pageInfo)
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

func (contentApi *ContentApi) Video(c *gin.Context) {

}

func (contentApi *ContentApi) Photo(c *gin.Context) {

}

func (contentApi *ContentApi) Search(c *gin.Context) {

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

}

func (contentApi *ContentApi) EditVideo(c *gin.Context) {

}

func (contentApi *ContentApi) EditPhoto(c *gin.Context) {

}

func (contentApi *ContentApi) List(c *gin.Context) {

}
