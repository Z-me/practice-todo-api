package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/handler"
	db "github.com/Z-me/practice-todo-api/middleware"
)

// Test test用
func Test() {
	nextID := db.GetNextID()
	fmt.Println("Open Todo API: ", nextID)
}

// Router main router
func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/todo", handler.GetTodoList)
	router.GET("/todo/:id", handler.GetTodoItemByID)
	router.POST("/todo", handler.AddNewTodo)
	router.PUT("/todo/:id", handler.UpdateTodoItem)
	router.PATCH("/todo/:id/status", handler.UpdateTodoState)
	router.DELETE("/todo/:id", handler.DeleteTodoListItem)

	return router
}
