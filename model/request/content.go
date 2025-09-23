package request

import "goMedia/model/appTypes"

type GetInfo struct {
	UID string `json:"uid" form:"uid"`
	PageInfo
}

type GetList struct {
	PageInfo
	TypeID  appTypes.TypeID `json:"type_id" form:"type_id"`
	Keyword string          `json:"keyword" form:"keyword"`
}

type ListByAdmin struct {
	PageInfo
	TypeID  int    `json:"type_id" form:"type_id"`
	Freeze  int    `json:"freeze" form:"freeze"`
	Keyword string `json:"keyword" form:"keyword"`
}

type GetID struct {
	UID string `json:"uid" form:"uid"`
}

type ContentFreeze struct {
	UID    string `json:"uid"`
	Freeze bool   `json:"freeze"`
}
