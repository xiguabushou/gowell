package service

type ServiceGroup struct {
	BaseService
	UserService
	JwtService
}

var ServiceGroupApp = new(ServiceGroup)
