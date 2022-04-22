package db

import (
	"fmt"

	"github.com/Z-me/practice-todo-api/api/model"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDb は確認用
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
	db.Create(&model.Todo{
		ID: 6,
		Title: "追加されたTodo",
		Status: "Backlog",
		Details: "自動追加のTodo",
		Priority: "P0",
	})

	// db.Find(&TodoList) // SELECT * FROM todos
	todoList := model.TodoList{}
	db.Find(&todoList)
	fmt.Println("todoList", todoList)

	// SELECT 2
	// db.Take(&todo) // SELECT * FROM todo LIMIT 1;
	// db.Last($todo) // SELECT * FROM todo ORDER BY id DESC LIMIT 1;
	todo := model.Todo{}
	db.Last(&todo)
	fmt.Println("Selected Title", todo.Title)

	// UPDATE
	db.Model(&todo).Update("Title", "更新されたTitle")
	fmt.Println("Updated Title", todo.Title)

	// DELETE
	db.Delete(&todo)
}

// GetNextID は次に指定するIDを取得する
func GetNextID() uint {
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return 0
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	todo := model.Todo{}
	db.Last(&todo)

	return todo.ID + 1
}

// GetTodoList DBからTodoリストを取得して返却する関数
func GetTodoList() model.TodoList {
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	todoList := model.TodoList{}
	db.Find(&todoList)

	return todoList
}

// AddNewTodo はDBに指定のPayloadの値を投入
func AddNewTodo(payload model.Payload) model.Todo {
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	newTodo := model.Todo{
		ID: GetNextID(),
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
	}

	db.Create(&newTodo)
	// db.Create(&model.Todo{
	// 	ID: GetNextID(),
	// 	Title: payload.Title,
	// 	Status: payload.Status,
	// 	Details: payload.Details,
	// 	Priority: payload.Priority,
	// })
	return newTodo
}
