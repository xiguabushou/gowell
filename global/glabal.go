package global

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"goMedia/config"
	"gorm.io/gorm"
)

var (
	Config     *config.Config
	Log        *zap.Logger
	BlackCache local_cache.Cache
	DB         *gorm.DB
)
