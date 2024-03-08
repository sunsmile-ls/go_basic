package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2. 生成唯一id
	userId := snowflake.GenID()
	user := models.ParamSignUp{
		UserId:     userId,
		Username:   p.Username,
		Password:   p.Password,
		RePassword: p.RePassword,
	}
	// 3. 插入数据库
	return mysql.InsertUser(&user)
}

func Login(p *models.ParamLogin) (string, error) {
	user := models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 判断用户是否存在
	if err := mysql.Login(&user); err != nil {
		return "", err
	}

	// 生成 token
	return jwt.GenToken(user.UserID, user.Username)
}
