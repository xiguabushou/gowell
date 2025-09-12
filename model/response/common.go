package response

type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type RecommendedList struct {
	Uid   string `json:"uid"`
	Cover string `json:"cover"`
	Title string `json:"title"`
}
