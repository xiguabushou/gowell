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

func (contentService *ContentService) GetList(info request.GetList) (any, int64, error) {
	db := global.DB

	if info.Keyword == "" {
		if info.TypeID == appTypes.VIDEO || info.TypeID == appTypes.PHOTO {
			db = db.Where("type_id = ? and freeze = ?", info.TypeID, appTypes.UnFreeze)
		} else {
			db = db.Where("freeze = ?", appTypes.UnFreeze)
		}
	} else {
		if info.TypeID == appTypes.VIDEO || info.TypeID == appTypes.PHOTO {
			db = db.Where("type_id = ? and freeze = ? and (title like ? or tags like ?)", info.TypeID, appTypes.UnFreeze, "%"+info.Keyword+"%", "%"+info.Keyword+"%")
		} else {
			db = db.Where("freeze = ? and (title like ? or tags like ?)", appTypes.UnFreeze, "%"+info.Keyword+"%", "%"+info.Keyword+"%")
		}
	}

	var pageinfo = request.PageInfo{
		Page:     info.Page,
		PageSize: info.PageSize,
	}

	option := other.MySQLOption{
		PageInfo: pageinfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Content{}, option)
}

func (contentService *ContentService) GetInfo(uid string) (response.GetInfo, error) {
	var content database.Content
	if err := global.DB.Where("uid = ? and freeze = ?", uid, appTypes.UnFreeze).First(&content).Error; err != nil {
		return response.GetInfo{}, err
	}

	var contentList []response.RecommendedList
	sql := `
        SELECT 
            *,
            (
                SELECT COUNT(*)
                FROM JSON_TABLE(t1.tags, '$[*]' COLUMNS (kw VARCHAR(50) PATH '$')) AS jt1
                WHERE JSON_CONTAINS(
                    (SELECT tags FROM contents WHERE uid = ?),
                    JSON_QUOTE(jt1.kw)
                )
            ) AS match_count
        FROM contents t1
        WHERE t1.uid != ? and t1.freeze = ?
        ORDER BY match_count DESC, t1.id
        LIMIT 6`
	if err := global.DB.Raw(sql, uid, uid, appTypes.UnFreeze).Scan(&contentList).Error; err != nil {
		return response.GetInfo{}, err
	}

	tags, err := utils.UnencodeJson(content.Tags)
	if err != nil {
		return response.GetInfo{}, nil
	}

	var resoult = response.GetInfo{
		Title:           content.Title,
		Video:           content.UID,
		Tags:            tags,
		RecommendedList: contentList,
	}
	return resoult, nil
}

func (contentService *ContentService) UploadVideo(title string, tags string, file *multipart.FileHeader, cover *multipart.FileHeader, c *gin.Context) error {
	NewUUID := uuid.Must(uuid.NewV4()).String()
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		unionTags, err := utils.EncodeJson(tags)
		if err != nil {
			return err
		}
		var newContent = database.Content{
			UID:    NewUUID,
			TypeID: appTypes.VIDEO,
			Title:  title,
			Tags:   unionTags,
			Freeze: appTypes.UnFreeze,
		}

		if err := c.SaveUploadedFile(cover, "uploads/video/"+NewUUID+"/cover.png"); err != nil {
			return errors.New("failed to save uploaded file")
		}
		if err := c.SaveUploadedFile(file, "uploads/video/"+NewUUID+"/video.mp4"); err != nil {
			return errors.New("failed to save uploaded file")
		}

		if err := tx.Create(&newContent).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		os.RemoveAll("uploads/video/" + NewUUID)
		return err
	}

	return nil
}

func (contentService *ContentService) UploadPhoto(title string, tags string, files []*multipart.FileHeader, cover *multipart.FileHeader, c *gin.Context) error {
	NewUUID := uuid.Must(uuid.NewV4()).String()
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		unionTags, err := utils.EncodeJson(tags)
		if err != nil {
			return err
		}
		num := 0

		if err := c.SaveUploadedFile(cover, "uploads/photo/"+NewUUID+"/cover.png"); err != nil {
			return err
		}
		for _, v := range files {
			photoID := uuid.Must(uuid.NewV4()).String()
			if err := c.SaveUploadedFile(v, "uploads/photo/"+NewUUID+"/"+photoID+".png"); err != nil {
				return err
			}
			num++
			var newPhoto = database.Photo{
				UID:     NewUUID,
				ImageID: photoID,
			}
			err := tx.Create(&newPhoto).Error
			if err != nil {
				return err
			}

		}

		var newContent = database.Content{
			UID:    NewUUID,
			TypeID: appTypes.PHOTO,
			Title:  title,
			Tags:   unionTags,
			Number: num,
			Freeze: appTypes.UnFreeze,
		}
		err = tx.Create(&newContent).Error
		return err
	})
	if err != nil {
		os.RemoveAll("uploads/photo" + NewUUID)
		return err
	}
	return nil
}

func (contentService *ContentService) ListByAdmin(info request.ListByAdmin) (any, int64, error) {
	db := global.DB

	if info.TypeID == appTypes.VIDEO || info.TypeID == appTypes.PHOTO {
		db = db.Where("type_id = ? and freeze = ?", info.TypeID, info.Freeze)
	}
	var pageinfo = request.PageInfo{
		Page:     info.Page,
		PageSize: info.PageSize,
	}

	option := other.MySQLOption{
		PageInfo: pageinfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Content{}, option)
}

func (contentService *ContentService) Freeze(uid string) error {
	var content database.Content
	return global.DB.Where("uid = ?", uid).First(&content).Update("freeze", appTypes.Freeze).Error
}

func (contentService *ContentService) UnFreeze(uid string) error {
	var content database.Content
	return global.DB.Where("uid = ?", uid).First(&content).Update("freeze", appTypes.UnFreeze).Error
}

func (contentService *ContentService) Delete(uid string) error {
	return  global.DB.Transaction( func(tx *gorm.DB) error {
		var content database.Content
		err := tx.Where("uid = ?",uid).First(&content).Error
		if err != nil {
			return err
		}

		if content.TypeID == appTypes.PHOTO{
			os.RemoveAll("uploads/photo/" + uid)
		}
		if content.TypeID == appTypes.VIDEO{
			os.RemoveAll("uploads/video/" + uid)
		}
		return nil
	})
}
