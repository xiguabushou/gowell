package service

import (
	"goMedia/global"
	"goMedia/model/database"
	"goMedia/model/other"
	"goMedia/model/request"
	"goMedia/utils"
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

func (contentService *ContentService) UploadVideo() {

}

func (contentService *ContentService) UploadPhoto() {

}

func (contentService *ContentService) EditVideo() {

}

func (contentService *ContentService) EditPhoto() {

}

func (contentService *ContentService) List() {

}
