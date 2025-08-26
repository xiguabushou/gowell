package flag

import (
	"errors"
	"fmt"
	"goMedia/global"
	"goMedia/model/appTypes"
	"goMedia/model/database"
	"goMedia/utils"

	"github.com/gofrs/uuid"
)

// Admin 用于创建一个管理员用户
func Admin() error {
	var user database.User

	// 提示用户输入邮箱
	fmt.Println("Enter email: ")
	// 读取用户输入的邮箱
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}

	// 提示用户输入密码
	fmt.Println("Enter password: ")
	// 读取用户输入密码
	var password string
	_, err = fmt.Scanln(&password)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Println("Enter password again: ")
	// 读取用户输入密码
	var password2 string
	_, err = fmt.Scanln(&password2)
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	if password != password2 {
		return fmt.Errorf("passwords do not match")
	}
	if err := checkPassword(password); err != nil {
		return err
	}

	// TODO 检测用户是否已存在

	NewUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to create UUID: %w", err)
	}

	// 填充用户数据
	user.UUID = NewUUID.String()
	user.Email = email
	user.Password = utils.BcryptHash(password)
	user.RoleID = appTypes.Admin
	user.Freeze = appTypes.UnFreeze

	// 在数据库中创建管理员用户
	if err := global.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("创建管理员失败: %w", err)
	}
	return nil
}

func checkPassword(password string) error {
	// 检查密码长度是否符合要求
	if len(password) < 6 || len(password) > 20 {
		return errors.New("password length should be between 6 and 20 characters")
	}
	return nil
}
