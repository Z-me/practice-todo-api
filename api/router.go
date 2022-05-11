package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/handler"
	"github.com/Z-me/practice-todo-api/middleware"
)

// Test test用
func Test() {
	fmt.Println("Open Todo API")
}

// Router main router
func Router() *gin.Engine {
	router := gin.Default()

	// router.GET("/todo", handler.GetTodoList)
	// router.GET("/todo/:id", handler.GetTodoItemByID)
	// router.POST("/todo", handler.AddNewTodo)
	// router.PUT("/todo/:id", handler.UpdateTodoItem)
	// router.PATCH("/todo/:id/status", handler.UpdateTodoState)
	// router.DELETE("/todo/:id", handler.DeleteTodoListItem)

	todoRouter := router.Group("/todo")
	todoRouter.Use(middleware.LoginCheckMiddleware())
	{
		todoRouter.GET("", handler.GetTodoList)
		todoRouter.GET("/todo/:id", handler.GetTodoItemByID)
		todoRouter.POST("/todo", handler.AddNewTodo)
		todoRouter.PUT("/todo/:id", handler.UpdateTodoItem)
		todoRouter.PATCH("/todo/:id/status", handler.UpdateTodoState)
		todoRouter.DELETE("/todo/:id", handler.DeleteTodoListItem)
	}

	return router
}
