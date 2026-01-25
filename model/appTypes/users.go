package appTypes

// RoleID 用户角色
type RoleID int

const (
	User  RoleID = iota //用户
	Vip                 //vip会员
	Admin               //管理员
)

type FreezeID bool

const (
	UnFreeze = false //未冻结
	Freeze   = true  //冻结
)
