package api

import "goMedia/service"

type ApiGroup struct {
	BaseApi
	UserApi
	ContentApi
}

var ApiGroupApp = new(ApiGroup)
var baseService = service.ServiceGroupApp.BaseService
var userService = service.ServiceGroupApp.UserService
var contentService = service.ServiceGroupApp.ContentService
