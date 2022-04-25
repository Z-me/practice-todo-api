package db

import (
	"fmt"
	"net/http"

	"github.com/Z-me/practice-todo-api/api/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*/ ConnectDb は確認用
func ConnectDb() {
	fmt.Println("start conect DB")
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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
*/

var dsn string = "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
var db *gorm.DB
var err error

// ConnectDB データベース接続
func ConnectDB() error {
	fmt.Println("Connect Database")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

// DisconnectDB データベースの接続解除
func DisconnectDB() {
	handleDb, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	defer handleDb.Close()
	fmt.Println("Dis-Connect Database")
}

// GetNextID は次に指定するIDを取得する
func GetNextID() uint {
	todo := model.Todo{}
	db.Last(&todo)
	return todo.ID + 1
}

// GetTodoList DBからTodoリストを取得して返却する関数
func GetTodoList() (model.TodoList, error) {
	todoList := model.TodoList{}
	err := db.Find(&todoList).Error
	return todoList, err
}

// GetTodoItemByID はIDをもとにItemを取得する関数
func GetTodoItemByID(id uint) (model.Todo, error) {
	todo := model.Todo{}
	err := db.Find(&todo, "id = ?", id).Error
	return todo, err
}

// AddNewTodo はDBに指定のPayloadの値を投入
func AddNewTodo(c *gin.Context, payload model.Payload) model.Todo {
	// dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to connect database"})
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

	if err := db.Create(&newTodo).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to create new item"})
	}
	return newTodo
}

// UpdateItem はDB上から指定のItemの情報を更新
func UpdateItem(c *gin.Context, id uint, payload model.Payload) model.Todo {
	// dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to connect database"})
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	target := model.Todo{}
	if err := db.Find(&target, "id = ?", id).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Todo List Item not found"})
	}

	updated := model.Todo{
		ID: id,
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
	}

	if err := db.Model(&target).Updates(&updated).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to update item"})
	}
	return updated
}

// UpdateItemStatus はDB上から指定のItemのStatusを更新
func UpdateItemStatus(c *gin.Context, id uint, status model.Status) model.Todo {
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to connect database"})
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	target := model.Todo{}
	if err := db.Find(&target, "id = ?", id).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Todo List Item not found"})
	}

	if err := db.Model(&target).Update("Status", status.Status).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to update item"})
	}

	return model.Todo{
		ID: id,
		Title: target.Title,
		Status: target.Status,
		Details: target.Details,
		Priority: target.Priority,
	}
}

// DeleteItem は任意のItemを削除
func DeleteItem(c *gin.Context, id uint) model.Todo {
	dsn := "host=localhost user=hajime.saito dbname=todo_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to connect database"})
	}
	handleDb, err := db.DB()
	defer handleDb.Close()

	target := model.Todo{}
	if err := db.Find(&target, "id = ?", id).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Todo List Item not found"})
	}

	result := model.Todo{
		ID: id,
		Title: target.Title,
		Status: target.Status,
		Details: target.Details,
		Priority: target.Priority,
	}

	if err := db.Delete(&target).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to update item"})
	}

	return result
}
