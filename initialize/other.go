package initialize

import (
	"go.uber.org/zap"
	"goMedia/global"
	"goMedia/utils"
	"os"

	"github.com/songzhibin97/gkit/cache/local_cache"
)

// OtherInit 执行其他配置初始化
func OtherInit() {
	// 解析令牌过期时间
	TokenExpiry, err := utils.ParseDuration(global.Config.Jwt.TokenExpiryTime)
	if err != nil {
		global.Log.Error("failed to parse refresh token expiry time configuration:", zap.Error(err))
		os.Exit(1)
	}

	// 配置本地缓存过期时间（使用刷新令牌过期时间，方便在远程登录或账户冻结时对 JWT 进行黑名单处理）
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(TokenExpiry),
	)
}
