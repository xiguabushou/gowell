package response

type GetInfo struct {
	Video           string            `json:"video"`
	Title           string            `json:"title"`
	Tags            []string          `json:"tags"`
	RecommendedList []RecommendedList `json:"recommend_list"`
	Images          []string          `json:"images"`
	Total           int               `json:"total"`
}
