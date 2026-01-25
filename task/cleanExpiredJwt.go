package task

import (
	"goMedia/global"
	"goMedia/model/database"
	"strconv"
	"time"
)

func CleanUpExpiredJwt() error{

	expiredTime := time.Now().Add(-7 * 24 * time.Hour)
	batchSize := 100 // 每批删除 100 条

	for {
		result := global.DB.Where("created_time < ?", expiredTime).Limit(batchSize).Delete(&database.JwtBlacklist{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			break // 没有更多数据可删除
		}

		msg := "Deleted batch of " + strconv.FormatInt(result.RowsAffected, 10) + " expired jwt"
		global.Log.Info(msg)

		// 可选：短暂休眠，减轻数据库压力
		time.Sleep(100 * time.Millisecond)

	}
	return nil
}
