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

	list, total, err := utils.MySQLPagination(&database.Content{}, option)
	if err != nil {
		return nil, 0, err
	}
	var cover string
	var contentType string
	var contentList []other.ContentList
	for _, v := range list {
		if v.TypeID == appTypes.VIDEO {
			cover = global.Config.System.Ip + "/video/" + v.UID + "/cover.png"
			contentType = "视频"
		}
		if v.TypeID == appTypes.PHOTO {
			cover = global.Config.System.Ip + "/photo/" + v.UID + "/cover.png"
			contentType = "图片"
		}
		content := other.ContentList{
			Uid:         v.UID,
			Title:       v.Title,
			Cover:       cover,
			ContentType: contentType,
		}
		contentList = append(contentList, content)
	}
	return contentList, total, err
}

func (contentService *ContentService) GetInfo(uid string, page int, pagesize int) (response.GetInfo, error) {
	var content database.Content
	if err := global.DB.Where("uid = ? and freeze = ?", uid, appTypes.UnFreeze).First(&content).Error; err != nil {
		return response.GetInfo{}, err
	}

	var contentList []database.Content
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
        WHERE t1.uid != ? and t1.freeze = ? and type_id = ?
        ORDER BY match_count DESC, t1.id
        LIMIT 6`
	if err := global.DB.Raw(sql, uid, uid, appTypes.UnFreeze, content.TypeID).Scan(&contentList).Error; err != nil {
		return response.GetInfo{}, err
	}

	tags, err := utils.UnencodeJson(content.Tags)
	if err != nil {
		return response.GetInfo{}, nil
	}

	if content.TypeID == appTypes.VIDEO {
		videoUrl := global.Config.System.Ip + "/video/" + content.UID + "/video.mp4"

		var newContentList []response.RecommendedList
		for _, v := range contentList {
			tempContent := response.RecommendedList{
				Uid:   v.UID,
				Cover: global.Config.System.Ip + "/video/" + v.UID + "/cover.png",
				Title: v.Title,
			}
			newContentList = append(newContentList, tempContent)
		}

		var resoult = response.GetInfo{
			Title:           content.Title,
			Video:           videoUrl,
			Tags:            tags,
			RecommendedList: newContentList,
		}
		return resoult, nil
	}

	if content.TypeID == appTypes.PHOTO {
		var newContentList []response.RecommendedList
		for _, v := range contentList {
			tempContent := response.RecommendedList{
				Uid:   v.UID,
				Cover: global.Config.System.Ip + "/photo/" + v.UID + "/cover.png",
				Title: v.Title,
			}
			newContentList = append(newContentList, tempContent)
		}

		var imagesUrl []string
		db := global.DB
		db = db.Where("uid = ?", uid)

		var pageinfo = request.PageInfo{
			Page:     page,
			PageSize: pagesize,
		}

		option := other.MySQLOption{
			PageInfo: pageinfo,
			Where:    db,
		}

		tempList, total, err := utils.MySQLPagination(&database.Photo{}, option)

		for _, v := range tempList {
			imageUrl := global.Config.System.Ip + "/photo/" + content.UID + "/" + v.ImageID + ".png"
			imagesUrl = append(imagesUrl, imageUrl)
		}

		var resoult = response.GetInfo{
			Title:           content.Title,
			Tags:            tags,
			RecommendedList: newContentList,
			Images:          imagesUrl,
			Total:           int(total),
		}
		return resoult, err
	}

	return response.GetInfo{}, nil
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

	if info.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+info.Keyword+"%")
	}

	if info.TypeID == 1 || info.TypeID == 2 {
		db = db.Where("type_id = ?", info.TypeID)
	}

	if info.Freeze == 0 || info.Freeze == 1 {
		db = db.Where("freeze = ?", info.Freeze)
	}

	var pageinfo = request.PageInfo{
		Page:     info.Page,
		PageSize: info.PageSize,
	}

	option := other.MySQLOption{
		PageInfo: pageinfo,
		Where:    db,
	}

	tempList, total, err := utils.MySQLPagination(&database.Content{}, option)
	if err != nil {
		return nil, 0, err
	}
	var cover string
	var contentList []other.ListByAdmin
	for _, v := range tempList {
		if v.TypeID == appTypes.VIDEO {
			cover = global.Config.System.Ip + "/video/" + v.UID + "/cover.png"
		}
		if v.TypeID == appTypes.PHOTO {
			cover = global.Config.System.Ip + "/photo/" + v.UID + "/cover.png"
		}
		content := other.ListByAdmin{
			Content:  v,
			Cover:     cover,
		}
		contentList = append(contentList, content)
	}
	return contentList, total, err

}

func (contentService *ContentService) Freeze(req request.ContentFreeze) error {
	var content database.Content
	return global.DB.Where("uid = ?", req.UID).First(&content).Update("freeze", req.Freeze).Error
}

func (contentService *ContentService) UnFreeze(uid string) error {
	var content database.Content
	return global.DB.Where("uid = ?", uid).First(&content).Update("freeze", appTypes.UnFreeze).Error
}

func (contentService *ContentService) Delete(uid string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var content database.Content
		err := tx.Where("uid = ?", uid).First(&content).Error
		if err != nil {
			return err
		}

		if content.TypeID == appTypes.PHOTO {
			os.RemoveAll("uploads/photo/" + uid)
		}
		if content.TypeID == appTypes.VIDEO {
			os.RemoveAll("uploads/video/" + uid)
		}
		return global.DB.Where("uid = ?", uid).Delete(&database.Content{}).Error
	})
}
