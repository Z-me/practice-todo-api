package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/handler"
)

func Test() {
	fmt.Println("Open Todo API")
}

func Router() *gin.Engine {
	router := gin.Default()

	handler.SetDefault()
	router.GET("/todo", handler.GetTodoList)
	router.GET("/todo/:id", handler.GetTodoItemById)
	router.POST("/todo", handler.AddNewTodo)
	router.POST("/todo/:id", handler.UpdateTodoItem)
	router.POST("/todo/:id/status/:status", handler.UpdateTodoState)
	router.DELETE("/todo/:id", handler.DeleteTodoListItem)

	return router
}
