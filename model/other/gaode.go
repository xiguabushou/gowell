package other

// IPResponse 用于表示 IP 定位查询的响应结果
type IPResponse struct {
	Status    string `json:"status"`    // 返回结果状态值：0表示失败，1表示成功
	Info      string `json:"info"`      // 返回状态说明
	InfoCode  string `json:"infocode"`  // 状态码：10000代表正确
	Province  string `json:"province"`  // 省份名称
	City      string `json:"city"`      // 城市名称
	Adcode    string `json:"adcode"`    // 城市的 adcode 编码
	Rectangle string `json:"rectangle"` // 所在城市矩形区域范围
}
