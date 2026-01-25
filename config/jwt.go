package config

type Jwt struct {
	TokenSecret     string `json:"token_secret" yaml:"token_secret"`           // 用于生成和验证访问令牌的密钥
	TokenExpiryTime string `json:"token_expiry_time" yaml:"token_expiry_time"` // 令牌的过期时间，例如 "15m" 表示 15 分钟
	Issuer          string `json:"issuer" yaml:"issuer"`                       // JWT 的签发者信息，通常是应用或服务的名称
}
