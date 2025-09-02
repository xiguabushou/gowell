package service

import (
	"errors"
	"goMedia/global"
	"goMedia/model/appTypes"
	"goMedia/model/database"
	"goMedia/model/other"
	"goMedia/model/request"
	"goMedia/model/response"
	"goMedia/utils"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
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

func (contentService *ContentService) Video(uid string) (resoult response.GetVideo,err error){
	var content database.Content
	if err := global.DB.Where("uid = ? ", uid).First(&content).Error; err != nil {
		return response.GetVideo{}, err
	}

	var contentList []response.RecommendedList
	sql := `
        SELECT 
            t1.name,
            t1.keywords,
            (
                SELECT COUNT(*)
                FROM JSON_TABLE(t1.keywords, '$[*]' COLUMNS (kw VARCHAR(50) PATH '$')) AS jt1
                WHERE JSON_CONTAINS(
                    (SELECT keywords FROM my_table WHERE name = ?),
                    JSON_QUOTE(jt1.kw)
                )
            ) AS match_count
        FROM my_table t1
        WHERE t1.name != ?
        ORDER BY match_count DESC, t1.id
        LIMIT 6`
	if err := global.DB.Raw(sql,uid,uid).Scan(&contentList).Error; err != nil{
		return response.GetVideo{} ,err
	}

	resoult = response.GetVideo{
		Title: content.Title,
		Video: content.UID,
		RecommendedList: contentList,
	}
	return resoult,nil
}

func (contentService *ContentService) Photo() {

}

func (contentService *ContentService) Search() {

}

func (contentService *ContentService) UploadVideo(title string, tags string, file *multipart.FileHeader, cover *multipart.FileHeader, c *gin.Context) error {
	NewUUID := uuid.Must(uuid.NewV4()).String()
	err := global.DB.Transaction(func(tx *gorm.DB) error{
		unionTags, err := utils.EncodeJson(tags)
		if err !=nil {
			return err
		}
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
			return err
		}
		return nil
	})
	if err != nil {
		os.Remove("uploads/video/" + NewUUID)
		return err
	}

	return nil
}

func (contentService *ContentService) UploadPhoto(title string, tags string, files []*multipart.FileHeader, cover *multipart.FileHeader, c *gin.Context) error {
	NewUUID := uuid.Must(uuid.NewV4()).String()
	err := global.DB.Transaction(func (tx *gorm.DB)error{
		unionTags, err:= utils.EncodeJson(tags)
		if err!=nil {
			return err
		}
		num := 0

		if err:= c.SaveUploadedFile(cover,"uploads/photo/" + NewUUID + "/cover.png");err != nil {
			return err
		}
		for _, v := range files {
			photoID := uuid.Must(uuid.NewV4()).String()
			if err:= c.SaveUploadedFile(v,"uploads/photo/" + NewUUID + "/" + photoID + ".png");err != nil{
				return err
			}
			num ++
			var newPhoto = database.Photo{
				UID: NewUUID,
				ImageID: photoID,
			}
		 	err := global.DB.Create(&newPhoto).Error
			if err!= nil{
				return err
			}
		
		}

		var newContent = database.Content{
			UID: NewUUID,
			TypeID: appTypes.PHOTO,
			Title: title,
			Tags: unionTags,
			Number: num,
		}
		err = global.DB.Create(&newContent).Error
		return err
	})
	if err != nil {
		os.Remove("uploads/photo" + NewUUID)
		return err
	}
	return nil
}

func (contentService *ContentService) EditVideo() {

}

func (contentService *ContentService) EditPhoto() {

}

func (contentService *ContentService) List() {

}
