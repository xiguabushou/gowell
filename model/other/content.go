package other

import "goMedia/model/database"

type ContentList struct {
	Uid         string `json:"uid"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	ContentType string `json:"content_type"`
}

type ListByAdmin struct {
	database.Content
	Cover string `json:"cover"`
}
