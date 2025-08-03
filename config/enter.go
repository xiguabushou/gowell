package config

type Config struct {
	Mysql   Mysql   `json:"mysql" yaml:"mysql"`
	Zap     Zap     `json:"zap" yaml:"zap"`
	Jwt     Jwt     `json:"jwt" yaml:"jwt"`
	System  System  `json:"system" yaml:"system"`
	Upload  Upload  `json:"upload" yaml:"upload"`
	Captcha Captcha `json:"captcha" yaml:"captcha"`
	Email   Email   `json:"email" yaml:"email"`
}
