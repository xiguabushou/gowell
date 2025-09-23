package api

import (
	"errors"
	"fmt"
	"goMedia/global"
	"goMedia/model/appTypes"
	"goMedia/model/database"
	"goMedia/model/request"
	"goMedia/model/response"
	"goMedia/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserApi struct{}

func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	session := sessions.Default(c)

	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	fmt.Println(savedEmail)
	fmt.Println(req.Email)
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("两次邮箱不一致！", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("验证码错误！", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("邮箱验证码过期！", c)
		return
	}

	NewUUID := uuid.Must(uuid.NewV4()).String()

	u := database.User{
		UUID:     NewUUID,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   appTypes.User,
	}

	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("账号已被注册！", c)
		return
	}

	// 注册成功后，生成 token 并返回
	userApi.TokenNext(c, user)
}

func (userApi *UserApi) Login(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("填写信息格式错误！", c)
		return
	}

	// 校验验证码
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := userService.Login(u)
		if err != nil {
			global.Log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
			return
		}

		// 登录成功后生成 token
		userApi.TokenNext(c, user)
		return
	}
	response.FailWithMessage("验证码错误！", c)
}

func (userApi *UserApi) AddUser(c *gin.Context) {
	var req request.AddUser
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage("数据输入格式错误或者长度不符合要求", c)
		return
	}
	err = userService.AddUser(req)
	if err != nil {
		global.Log.Error("Failed to add user:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("User added successfully", c)
}

func (userApi *UserApi) DeleteUser(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage("数据绑定失败,请检查格式与长度", c)
		return
	}
	err = userService.DeleteUser(req.ID)
	if err != nil {
		global.Log.Error("Failed to delete user:", zap.Error(err))
		response.FailWithMessage("删除用户失败", c)
		return
	}
	response.OkWithMessage("User deleted successfully", c)
}

func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("数据绑定失败,请检查格式与长度", c)
		return
	}

	if store.Verify(req.CaptchaID, req.Captcha, true) {
		err = userService.ForgotPassword(req.Email)
		if err != nil {
			global.Log.Error("Failed to send email:", zap.Error(err))
			response.FailWithMessage("发送邮件失败", c)
			return
		}
		response.OkWithMessage("Successfully sent email", c)
		return
	}
	response.FailWithMessage("无效的验证码", c)
}

func (userApi *UserApi) AskForVip(c *gin.Context) {
	var req request.AskForVip
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage("数据绑定失败,请检查格式与长度", c)
		return
	}
	err = userService.AskForVip(req)
	if err != nil {
		global.Log.Error("Application failed:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("The application was successful", c)
}

func (userApi *UserApi) Logout(c *gin.Context) {
	if err := userService.Logout(c); err != nil {
		global.Log.Error("Failed to logout:", zap.Error(err))
		response.FailWithMessage("注销失败", c)
		return
	}
	response.OkWithMessage("Successful logout", c)
}

func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("数据绑定失败,请检查格式与长度", c)
		return
	}
	req.UUID = utils.GetUUID(c)
	err = userService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("Failed to modify:", zap.Error(err))
		response.FailWithMessage("修改密码失败,密码错误", c)
		return
	}
	response.OkWithMessage("Successfully changed password, please log in again", c)
}

func (userApi *UserApi) UserInfo(c *gin.Context) {
	UUID := utils.GetUUID(c)
	user, err := userService.UserInfo(UUID)
	if err != nil {
		global.Log.Error("Failed to get user information:", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}
	response.OkWithData(user, c)
}

func (userApi *UserApi) EditUser(c *gin.Context) {
	var req request.EditUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.Log.Error("Failed to get edit information:", zap.Error(err))
		response.FailWithMessage("Failed to get edit information", c)
		return
	}
	err = userService.EditUser(req)
	if err != nil {
		global.Log.Error("Failed to edit user:", zap.Error(err))
		response.FailWithMessage("Failed to edit user", c)
		return
	}
	response.OkWithMessage("Successfully edit user", c)
}

func (userApi *UserApi) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user login list:", zap.Error(err))
		response.FailWithMessage("Failed to get user login list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (UserApi *UserApi) GetListAboutAskForVip(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.GetListAboutAskForVip(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get ask for vip list:", zap.Error(err))
		response.FailWithMessage("Failed to get ask for vip list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

func (userApi *UserApi) ResetForgotPassword(c *gin.Context) {
	var req request.ResetForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetForgotPassword(req)
	if err != nil {
		global.Log.Error("Failed to reset forgot password:", zap.Error(err))
		response.FailWithMessage("Failed to reset forgot password", c)
		return
	}
	response.OkWithMessage("Successfully reset forgot password", c)
}

func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	// 检查用户是否被冻结
	if user.Freeze {
		response.FailWithMessage("The user is frozen, contact the administrator", c)
		return
	}

	baseClaims := request.BaseClaims{
		Email: user.Email,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()

	// 创建访问令牌
	accessClaims := j.CreateAccessClaims(baseClaims)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get accessToken:", zap.Error(err))
		response.FailWithMessage("Failed to get accessToken", c)
		return
	}

	// 是否开启了多地点登录拦截
	// 未开启多地点登录拦截
	if !global.Config.System.UseMultipoint {
		c.Set("user_id", user.UUID)
		c.Set("email", user.Email)
		response.OkWithDataAndMessage(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
		return
	}

	// 开启多地点登录拦截

	// 检查 mysql jwt 中是否已存在该用户的 JWT
	var oldJwt database.Jwt
	err = global.DB.Where("uuid = ?", user.UUID).First(&oldJwt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 无旧 token
			// 添加到jwt表
			newJwt := database.Jwt{
				UUID:        user.UUID,
				Jwt:         accessToken,
				UpdatedTime: time.Now(),
			}
			global.DB.Create(&newJwt)

			c.Set("user_id", user.UUID)
			c.Set("email", user.Email)
			response.OkWithDataAndMessage(response.Login{
				User:                 user,
				AccessToken:          accessToken,
				AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
			}, "Successful login", c)
			return
		} else {
			// 处理其他错误
			global.Log.Error("mysql error to find jwt ", zap.Error(err))
			response.FailWithMessage("mysql error to find jwt ", c)
			return
		}
	}
	dr, err := utils.ParseDuration(global.Config.Jwt.TokenExpiryTime)
	if err != nil {
		global.Log.Error("Failed to get Parse TokenExpiryTime:", zap.Error(err))
		response.FailWithMessage("Failed to get Parse TokenExpiryTime", c)
		return
	}
	// jwt未过期
	if oldJwt.UpdatedTime.Add(dr).Unix() > time.Now().Unix() {
		// 将当前jwt加入黑名单
		var jwtBlacklist database.JwtBlacklist
		jwtBlacklist.Jwt = oldJwt.Jwt
		jwtBlacklist.CreatedTime = oldJwt.UpdatedTime
		if err := global.DB.Create(&jwtBlacklist).Error; err != nil {
			global.Log.Error("Mysql error to create jwtBlacklist", zap.Error(err))
			response.FailWithMessage("Mysql error to create jwtBlacklist", c)
			return
		}
		global.BlackCache.SetDefault(oldJwt.Jwt, struct{}{})

		//更新jwt表
		oldJwt.UpdatedTime = time.Now()
		oldJwt.Jwt = accessToken
		if err := global.DB.Save(&oldJwt).Error; err != nil {
			global.Log.Error("Mysql save jwt blacklist error", zap.Error(err))
			response.FailWithMessage("Mysql save jwt blacklist error", c)
			return
		}

		c.Set("user_id", user.UUID)
		c.Set("email", user.Email)
		response.OkWithDataAndMessage(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
		return
	}
	//更新jwt表
	oldJwt.UpdatedTime = time.Now()
	oldJwt.Jwt = accessToken
	if err := global.DB.Save(&oldJwt).Error; err != nil {
		global.Log.Error("Mysql save jwt blacklist error", zap.Error(err))
		response.FailWithMessage("Mysql save jwt blacklist error", c)
		return
	}

	c.Set("user_id", user.UUID)
	c.Set("email", user.Email)
	response.OkWithDataAndMessage(response.Login{
		User:                 user,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
	}, "Successful login", c)
}

func (userApi *UserApi) ApprovingForVip(c *gin.Context) {
	var req request.ApprovingForVip
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ApprovingForVip(req, c)
	if err != nil {
		global.Log.Error("The commit failed:", zap.Error(err))
		response.FailWithMessage("The commit failed", c)
		return
	}
	response.OkWithMessage("The commit success", c)
}
