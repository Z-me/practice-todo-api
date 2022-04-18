package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	todoType "github.com/Z-me/practice-todo-api/api/model"
)

var todoList = []todoType.Todo{}

// NOTE: Default Todo
func SetDefault() {
	todoList = []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}
}

func GetTodoList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todoList)
}

func GetTodoListItemById(c *gin.Context) {
	id := c.Param("id")
	for _, item := range todoList {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}

func PostTodoItem(c *gin.Context) {
	var newTodo todoType.Todo

	if err := c.BindJSON(&newTodo); err != nil {
		return
	}

	todoList = append(todoList, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func DeleteTodoListItemById(c *gin.Context) {
	id := c.Param("id")
	for i, item := range todoList {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			todoList = append(todoList[:i], todoList[i + 1:]...)
			fmt.Println("Deleted Todolist", todoList)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}

func PostTodoListItemById(c *gin.Context) {
	id := c.Param("id")
	var newTodo todoType.Todo

	if err := c.BindJSON(&newTodo); err != nil {
		return
	}

	for i, item := range todoList {
		if item.ID == id {
			tmp := append(todoList[:i], newTodo)
			todoList = append(tmp, todoList[i + 1:]...)
			fmt.Println("Updated Todolist", newTodo)
			c.IndentedJSON(http.StatusCreated, newTodo)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
	// todoList = append(todoList, newTodo)
}

func PostTodoListItemUpdateStateById(c *gin.Context) {
	id := c.Param("id")
	status := c.Param("status")

	for i, item := range todoList {
		if item.ID == id {
			target := item
			target.Status = status
			tmp := append(todoList[:i], target)
			todoList = append(tmp, todoList[i + 1:]...)
			fmt.Println("Updated Todolist", target)
			c.IndentedJSON(http.StatusCreated, target)
			return
		}
	}

}
