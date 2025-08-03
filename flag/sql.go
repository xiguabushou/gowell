package flag

import (
	"goMedia/global"
	"goMedia/model/database"
)

// SQL 表结构迁移，如果表不存在，它会创建新表；如果表已经存在，它会根据结构更新表
func SQL() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&database.User{},
		&database.Login{},
		&database.JwtBlacklist{},
		&database.Jwt{},

		// TODO 添加更新的数据库表
	)
}
