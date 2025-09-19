package initialize

import (
	"goMedia/global"
	"goMedia/middleware"
	"goMedia/router"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()
	// 使用日志记录中间件
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	// 使用gin会话路由
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		Secure:   false, // ⚠️ 开发环境设为 false
		HttpOnly: true,
	})

	Router.Use(sessions.Sessions("session", store))
	// 将指定目录下的文件提供给客户端
	// "uploads" 是URL路径前缀，http.Dir("uploads")是实际文件系统中存储文件的目录
	Router.StaticFS(global.Config.Upload.Path, http.Dir(global.Config.Upload.Path))
	// 创建路由组
	routerGroup := router.RouterGroupApp

	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	userGroup := Router.Group(global.Config.System.RouterPrefix)
	//userGroup.Use(middleware.JWTAuth())
	vipGroup := Router.Group(global.Config.System.RouterPrefix)
	//vipGroup.Use(middleware.JWTAuth()).Use(middleware.VipAuth())
	adminGroup := Router.Group(global.Config.System.RouterPrefix)
	//adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
	{
		routerGroup.BaseRouter.InitBaseRouter(publicGroup)
		routerGroup.UserRouter.InitUserRouter(publicGroup, userGroup, vipGroup, adminGroup)
		routerGroup.ContentRouter.InitContentRouter(vipGroup, adminGroup)
	}

	return Router
}
