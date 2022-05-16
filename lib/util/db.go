package util

import (
	"fmt"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn string = "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
var db *gorm.DB
var err error

// ConnectDB データベース接続
func ConnectDB() error {
	fmt.Println("Connect Database")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

// DisconnectTodoDB データベースの接続解除
func DisconnectDB() {
	handleDb, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	defer handleDb.Close()
	fmt.Println("Dis-Connect Database")
}

// GetDbObj データベースObjectの取得
func GetDbObj() *gorm.DB {
	return db
}

// テスト用の データベース切り替え
func UseTestBD() {
	dsn = "host=localhost user=hajime.saito dbname=todo_app_test port=5432 sslmode=disable"
}
