package middleware

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/Z-me/practice-todo-api/api/model"
)

type Todo struct {
	ID        	uint		`gorm:"primaryKey"`
	Title     	string
	Status		string
	Details		string
	Priority	string
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type TodoList []Todo

func ConnectDb() {
	fmt.Println("start conect DB")
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	// INSERT
	// 一気に入れてみたパターン
	// db.Create(&[]Todo{
	// 	{ID: 1,	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
	// 	{ID: 2,	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
	// 	{ID: 3,	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
	// 	{ID: 4,	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
	// 	{ID: 5,	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	// })
	db.Create(&Todo{
		ID: 6,
		Title: "追加されたTodo",
		Status: "Backlog",
		Details: "自動追加のTodo",
		Priority: "P0",
	})

	// db.Find(&TodoList) // SELECT * FROM todos
	todoList := TodoList{}
	db.Find(&todoList)
	fmt.Println("todoList", todoList)

	// SELECT 2
	// db.Take(&todo) // SELECT * FROM todo LIMIT 1;
	// db.Last($todo) // SELECT * FROM todo ORDER BY id DESC LIMIT 1;
	todo := Todo{}
	db.Last(&todo)
	fmt.Println("Selected Title", todo.Title)

	// UPDATE
	db.Model(&todo).Update("Title", "更新されたTitle")
	fmt.Println("Updated Title", todo.Title)

	// DELETE
	db.Delete(&todo)
}
