package service

type ServiceGroup struct {
	BaseService
	UserService
	JwtService
	ContentService
}

var ServiceGroupApp = new(ServiceGroup)
