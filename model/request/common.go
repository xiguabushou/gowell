package request

import "goMedia/model/appTypes"

type PageInfo struct {
	Page     int             `json:"page" form:"page"`
	PageSize int             `json:"page_size" form:"page_size"`
	TypeID   appTypes.TypeID `json:"type_id" form:"type_id"`
}
