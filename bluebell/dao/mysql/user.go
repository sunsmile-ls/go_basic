package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 数据库的每个语句要原子性
const secret = "sunsmile.com"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

func CheckUserExist(username string) (err error) {
	// 查询数据库，判断用户是否存在
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 插入用户
func InsertUser(user *models.ParamSignUp) (err error) {
	password := encryptPassword(user.Password)
	sqlStr := "insert into user(user_id, username, password) values (?, ?, ?)"
	_, err = db.Exec(sqlStr, user.UserId, user.Username, password)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 根据用户名判断用户是否存在
func Login(user *models.User) (err error) {
	oldPassword := user.Password
	// 通过用户名获取密码
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	// 判断用户名和密码是否正确
	password := encryptPassword(oldPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return nil
}
