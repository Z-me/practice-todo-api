package db

import (
	"github.com/Z-me/practice-todo-api/api/model"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getUserID は認証に利用されたユーザーのIDを取得
func getUserID(dbObj *gorm.DB, name string) uint {
	user := model.User{}
	dbObj.Where("name = ?", name).First(&user)
	return user.ID
}

// CheckUserAuth は認証のmiddlewareで呼び出されるユーザー認証用関数
func CheckUserAuth(dbObj *gorm.DB, name, password string) bool {
	user := model.User{}
	if result := dbObj.Where("name = ?", name).First(&user); result.Error != nil {
		return false
	}
	if password != user.Password {
		return false
	}
	return true
}
