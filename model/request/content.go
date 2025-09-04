package request

import "goMedia/model/appTypes"

type GetInfo struct {
	UID string `json:"uid" form:"uid"`
}

type GetList struct {
	PageInfo
	TypeID  appTypes.TypeID `json:"type_id" form:"type_id"`
	Keyword string          `json:"keyword" form:"keyword"`
}

type ListByAdmin struct {
	PageInfo
	TypeID  appTypes.TypeID   `json:"type_id" form:"type_id"`
	Freeze  appTypes.FreezeID `json:"freeze" form:"freeze"`
	Keyword string            `json:"keyword" form:"keyword"`
}

type GetID struct {
	UID string `json:"uid" form:"uid"`
}

type EditTitleAndTags struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
	Tags  string `json:"tags"`
}

type DeleteContentVideo struct {
	UID string `json:"uid"`
	Name string `json:"name"`
}

type DeleteContentPhoto struct{
	UID string `json:"uid"`
	ImageID []string `json:"image_id"`
}