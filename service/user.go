package service

import (
	"errors"
	"goMedia/global"
	"goMedia/model/appTypes"
	"goMedia/model/database"
	"goMedia/model/other"
	"goMedia/model/request"
	"goMedia/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
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
			return database.User{}, errors.New("账号或密码错误")
		}
		return user, nil
	}

	//记录错误不存在的登录信息
	var filedLoginMsg = database.FailedLogin{
		Email:    u.Email,
		Password: u.Password,
	}
	global.DB.Create(&filedLoginMsg)
	return database.User{}, err
}

func (userService *UserService) AddUser(u request.AddUser) error {
	var user database.User
	if gorm.ErrRecordNotFound != global.DB.Where("email = ?", u.Email).First(&database.User{}).Error {
		return errors.New("该账户已存在")
	}
	user.Email = u.Email
	user.Password = utils.BcryptHash(u.Password)
	user.UUID = uuid.Must(uuid.NewV4()).String()
	user.Freeze = u.Freeze
	user.RoleID = u.RoleID
	if global.DB.Create(&user).Error != nil {
		return errors.New("添加用户失败")
	}
	return nil
}

func (userService *UserService) DeleteUser(id string) error {
	return global.DB.Where("uuid = ?", id).Delete(&database.User{}).Error
}

func (userService *UserService) ForgotPassword(email string) error {

	baseClaims := request.ForgotPasswordClaims{
		Email: email,
	}

	j := utils.NewJWT()
	TokenClaims := j.CreateTokenClaims(baseClaims)
	token, _ := j.CreateToken(TokenClaims)

	url := global.Config.System.ForgotPasswordUrl + "?token=" + token

	subject := "重新设置密码请求"
	body := `亲爱的用户[` + email + `]，
<br/>
<br/>
你正在进行使用该邮箱重新设置密码！为了确保您的账号安全，请进入下面网址进行修改密码：<br/>
<br/>
网址：[<font color="blue"><u>` + url + `</u></font>]<br/>
该网址在 5 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有此请求，请忽略此邮件。
<br/>
`
	_ = utils.Email(email, subject, body)
	return nil
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
	if err := global.DB.Where("uuid = ?", req.UUID).First(&user).Error; err != nil {
		return err
	}
	if req.Password != "" {
		user.Password = utils.BcryptHash(req.Password)
	}
	user.RoleID = req.RoleID
	user.Freeze = req.Freeze
	return global.DB.Save(&user).Error

}
func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB

	if info.Search != "" {
		db = db.Where("email LIKE ?", "%"+info.Search+"%")
	}

	if info.RoleID == 0 || info.RoleID == 1 || info.RoleID == 2 {
		db = db.Where("role_id = ?", info.RoleID)
	}

	if info.IsFreeze == 0 || info.IsFreeze == 1 {
		db = db.Where("freeze = ?", info.IsFreeze)
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

func (userService *UserService) ResetForgotPassword(req request.ResetForgotPassword) error {
	token := req.Token

	j := utils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		return err
	}

	var JwtBlacklist database.JwtBlacklist
	if !errors.Is(global.DB.Where("jwt = ?", token).First(&JwtBlacklist).Error, gorm.ErrRecordNotFound) {
		return errors.New("jwt blacklist")
	}

	var user database.User
	if err := global.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)

	oldJwt := database.JwtBlacklist{
		Jwt:         token,
		CreatedTime: time.Now(),
	}

	global.DB.Create(&oldJwt)
	return global.DB.Save(&user).Error
}

func (userService *UserService) AskForVip(req request.AskForVip) error {
	var oldAskForVip database.AskForVip

	err := global.DB.Where("uuid = ?", req.UUID).Where("finish_at is null").First(oldAskForVip).Error
	if err == nil {
		// 找到了记录，说明存在相同的未完成请求
		return errors.New("the same request already exists")
	}

	var newAskForVip = database.AskForVip{
		Email:    req.Email,
		Message:  req.Message,
		UUID:     req.UUID,
		FinishAt: nil,
	}

	return global.DB.Create(&newAskForVip).Error
}

func (userService *UserService) GetListAboutAskForVip(info request.PageInfo) (any, int64, error) {
	db := global.DB
	db = db.Where("finish_at is NULL")

	var options = other.MySQLOption{
		PageInfo: info,
		Where:    db,
	}

	return utils.MySQLPagination(&database.AskForVip{}, options)
}

func (userService *UserService) ApprovingForVip(req request.ApprovingForVip, c *gin.Context) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if req.IsPass {
			var user database.User
			err := tx.Where("uuid = ?", req.UUID).First(&user).Update("role_id", appTypes.Vip).Error
			if err != nil {
				return err
			}
			var ask database.AskForVip
			err = tx.Where("uuid = ? and finish_at is null", req.UUID).First(&ask).Error
			if err != nil {
				return err
			}
			t := time.Now()
			ask.FinishAt = &t
			ask.Approver = utils.GetEmail(c)
			ask.ApprovalResults = req.IsPass
			return tx.Save(&ask).Error
		} else {
			var ask database.AskForVip
			err := tx.Where("uuid = ? and finish_at is null", req.UUID).First(&ask).Error
			if err != nil {
				return err
			}
			t := time.Now()
			ask.FinishAt = &t
			ask.Approver = utils.GetEmail(c)
			ask.ApprovalResults = req.IsPass
			return tx.Save(&ask).Error
		}
	})
}
