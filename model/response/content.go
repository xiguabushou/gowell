package response

type GetVideo struct {
	Video           string            `json:"video"`
	Title           string            `json:"titile"`
	RecommendedList []RecommendedList `json:"recommend_list"`
}
