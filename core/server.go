package core

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goMedia/global"
	"goMedia/initialize"
	"net/http"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	addr := global.Config.System.Addr()
	Router := initialize.InitRouter()

	// TODO 加载黑名单

	// 初始化服务器并启动
	s := initServer(addr, Router)
	global.Log.Info("server run success on ", zap.String("address", addr))
	global.Log.Error(s.ListenAndServe().Error())
}

// initServer 函数初始化一个标准的 HTTP 服务器
func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,          // 设置服务器监听的地址
		Handler:        router,           // 设置请求处理器（路由）
		ReadTimeout:    10 * time.Minute, // 设置请求的读取超时时间为 10 分钟
		WriteTimeout:   10 * time.Minute, // 设置响应的写入超时时间为 10 分钟
		MaxHeaderBytes: 1 << 20,          // 设置最大请求头的大小（1MB）
	}
}
