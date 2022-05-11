package db

import (
	"fmt"

	"github.com/Z-me/practice-todo-api/api/model"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var userDsn string = "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
var userDb *gorm.DB

// ConnectUserDB user tableの接続
func ConnectUserDB() error {
	fmt.Println("Connect Database [users table]")
	userDb, err = gorm.Open(postgres.Open(userDsn), &gorm.Config{})
	return err
}

// DisconnectUserDB user tableの接続解除
func DisconnectUserDB() {
	handleDb, err := userDb.DB()
	if err != nil {
		panic("failed to connect database [users table]")
	}
	defer handleDb.Close()
	fmt.Println("Dis-Connect Database [users table]")
}

// getUserID は認証に利用されたユーザーのIDを取得
func getUserID(name string) uint {
	user := model.User{}
	userDb.Where("name = ?", name).First(&user)
	return user.ID
}

// CheckUserAuth は認証のmiddlewareで呼び出されるユーザー認証用関数
func CheckUserAuth(name, password string) bool {
	ConnectUserDB()
	defer DisconnectUserDB()
	user := model.User{}
	if result := userDb.Where("name = ?", name).First(&user); result.Error != nil {
		return false
	}
	if password != user.Password {
		return false
	}
	return true
}
