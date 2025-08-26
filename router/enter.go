package router

type RouterGroup struct {
	BaseRouter
	UserRouter
	ContentRouter
}

var RouterGroupApp = new(RouterGroup)
