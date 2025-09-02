package response

type GetInfo struct {
	Video           string            `json:"video"`
	Title           string            `json:"titile"`
	Tags            []string            `json:"tags"`
	RecommendedList []RecommendedList `json:"recommend_list"`
}
