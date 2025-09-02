package response


type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type RecommendedList struct {
	UID string `json:"uid"`
	Title string `json:"title"`
}
