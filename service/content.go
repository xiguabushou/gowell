package service

import (
	"errors"
	"goMedia/global"
	"goMedia/model/appTypes"
	"goMedia/model/database"
	"goMedia/model/other"
	"goMedia/model/request"
	"goMedia/utils"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type ContentService struct{}

func (contentService *ContentService) Home(info request.PageInfo) (interface{}, int64, error) {
	db := global.DB

	option := other.MySQLOption{
		PageInfo: info,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Content{}, option)
}

func (contentService *ContentService) Video() {

}

func (contentService *ContentService) Photo() {

}

func (contentService *ContentService) Search() {

}

func (contentService *ContentService) UploadVideo(title string, tags string, file *multipart.FileHeader, cover *multipart.FileHeader, c *gin.Context) error {
	NewUUID := uuid.Must(uuid.NewV4()).String()
	unionTags := strings.ReplaceAll(tags,"，",",")
	var newContent = database.Content{
		UID: NewUUID,
		TypeID: appTypes.VIDEO,
		Title: title,
		Tags: unionTags,
	}

 	if err := c.SaveUploadedFile(cover,"uploads/video/" + NewUUID + "/cover.png"); err != nil {
		return errors.New("failed to save uploaded file")
	}
	if err := c.SaveUploadedFile(file,"uploads/video/" + NewUUID + "/video.mp4"); err != nil {
		return errors.New("failed to save uploaded file")
	}

	if err := global.DB.Create(&newContent).Error; err != nil {
		os.Remove("uploads/video/" + NewUUID + "/video.mp4")
		os.Remove("uploads/video/" + NewUUID + "/cover.png")
		return err
	}

	return nil
}

func (contentService *ContentService) UploadPhoto() {

}

func (contentService *ContentService) EditVideo() {

}

func (contentService *ContentService) EditPhoto() {

}

func (contentService *ContentService) List() {

}
