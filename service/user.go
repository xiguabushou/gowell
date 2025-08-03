package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"goMedia/global"
	"goMedia/model/database"
	"goMedia/model/other"
	"goMedia/model/request"
	"goMedia/utils"
	"gorm.io/gorm"
	"time"
)

type UserService struct{}

func (userService *UserService) Register(u database.User) (user database.User, err error) {
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&user).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}

	// 密码加密
	u.Password = utils.BcryptHash(u.Password)
	// TODO 关闭加密注释上一行即可

	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}

	return u, nil
}

func (userService *UserService) Login(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err
}

func (userService *UserService) AddUser(u request.AddUser) error {
	var user database.User
	user.Email = u.Email
	user.Password = utils.BcryptHash(u.Password)
	user.UUID = uuid.Must(uuid.NewV4()).String()
	user.Freeze = u.Freeze
	user.RoleID = u.RoleID
	return global.DB.Create(&user).Error
}

func (userService *UserService) DeleteUser(id string) error {
	return global.DB.Where("uuid = ?", id).Delete(&database.User{}).Error
}

func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) Logout(c *gin.Context) error {
	uuid := utils.GetUUID(c)
	utils.ClearAccessToken(c)
	var jwt database.Jwt
	return global.DB.Where("uuid = ?", uuid).Delete(&jwt).Error
}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UUID).Error; err != nil {
		return err
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("original password does not match the current account")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}
func (userService *UserService) UserInfo(UUID string) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, UUID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}
func (userService *UserService) EditUser(req request.EditUser) error {
	var user database.User
	if err := global.DB.Take(&user, req.UUID).Error; err != nil {
		return err
	}
	if req.Password != "" {
		user.Password = utils.BcryptHash(req.Password)
	}
	user.Email = req.Email
	user.RoleID = req.RoleID
	user.Freeze = req.Freeze
	return global.DB.Save(&user).Error

}
func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB

	if info.Search != "" {
		db = db.Where("email LIKE ?", "%"+info.Search+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.User{}, option)

}
func (userService *UserService) UserFreeze(req request.UserOperation) error {
	var user database.User
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}

	var jwt database.Jwt
	if err := global.DB.Take(&jwt, user.ID).Error; err != nil {
		return err
	}

	if jwt.Jwt != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwt.Jwt, CreatedTime: time.Now()})
	}

	return nil
}
func (userService *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(&database.User{}, req.ID).Update("freeze", false).Error
}
func (userService *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {
	db := global.DB

	if info.Search != "" {
		db = db.Where("email LIKE ?", "%"+info.Search+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}

	return utils.MySQLPagination(&database.Login{}, option)
}
