package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/model"
)

var todoList = []model.Todo{}
// Note: Structのリテラルとしてmodel.Idだと使えないらしい
// var nextId model.Id
var nextId int

// NOTE: Default Todo
func LoadInitialData() {
	todoList = []model.Todo{
		{ID: 1,	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: 2,	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: 3,	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: 4,	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: 5,	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}
	nextId = 6
}

func GetTodoList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todoList)
}

func GetTodoItemById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, item := range todoList {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}

func AddNewTodo(c *gin.Context) {
	var payload model.Payload

	if err := c.BindJSON(&payload); err != nil {
		return
	}

	newTodo := model.Todo {
		ID: nextId,
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
	}
	nextId +=  1

	todoList = append(todoList, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func UpdateTodoItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload model.Payload

	if err := c.BindJSON(&payload); err != nil {
		return
	}

	newTodo := model.Todo {
		ID: id,
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
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
}

func UpdateTodoState(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload model.StatusPayload

	if err := c.BindJSON(&payload); err != nil {
		return
	}

	for i, item := range todoList {
		if item.ID == id {
			target := item
			target.Status = payload.Status
			tmp := append(todoList[:i], target)
			todoList = append(tmp, todoList[i + 1:]...)
			fmt.Println("Updated Todolist", target)
			c.IndentedJSON(http.StatusCreated, target)
			return
		}
	}
}

func DeleteTodoListItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, item := range todoList {
		if item.ID == id {
			todoList = append(todoList[:i], todoList[i + 1:]...)
			c.IndentedJSON(http.StatusOK, todoList)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}
