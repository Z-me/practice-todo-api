package todoApi

import (
	"fmt"

	"github.com/gin-gonic/gin"

	handler "github.com/Z-me/practice-todo-api/api/handler"
)

func Test() {
	fmt.Println("Open Todo API")
}

func Router() {
	router := gin.Default()
	router.GET("/todo", handler.GetTodoList)
	router.GET("/todo/:id", handler.GetTodoListItemById)
	router.POST("/todo", handler.PostTodoItem)
	router.POST("/todo/:id", handler.PostTodoListItemById)
	router.POST("/todo/:id/status/:status", handler.PostTodoListItemUpdateStateById)
	router.DELETE("/todo/:id", handler.DeleteTodoListItemById)

	router.Run("localhost:8080")
}
