package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goMedia/utils"
	"time"
)

type BaseService struct{}

func (baseService BaseService) SendEmailVerificationCode(c *gin.Context, to string) error {
	// 获取随机数字
	verificationCode := utils.GenerateVerificationCode(6)
	// 设置过期时间
	expireTime := time.Now().Add(5 * time.Minute).Unix()

	// 将验证码、验证邮箱、过期时间存入会话中
	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("expire_time", expireTime)
	session.Set("email", to)
	_ = session.Save()

	subject := "您的邮箱验证码"
	body := `亲爱的用户[` + to + `]，
<br/>
<br/>
你正在进行使用该邮箱进行注册！为了确保您的邮箱安全，请使用以下验证码进行验证：<br/>
<br/>
验证码：[<font color="blue"><u>` + verificationCode + `</u></font>]<br/>
该验证码在 5 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有请求此验证码，请忽略此邮件。
<br/>
`
	_ = utils.Email(to, subject, body)
	return nil
}
