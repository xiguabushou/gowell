package database

type Photo struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	UID     string `json:"uid"`
	ImageID string `json:"image_id"`
}
