package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/model"
	db "github.com/Z-me/practice-todo-api/middleware"
)

type Status string

type Todo struct {
  ID        int     `json:"id"`
  Title     string  `json:"title" binding:"required,max=30"`
  Status    string  `json:"status" binding:"required"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority" binding:"required,max=1000"`
}

type Payload struct {
  Title     string  `json:"title" binding:"required,max=30"`
  Status    string  `json:"status" binding:"required"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority" binding:"required,max=1000"`
}

type StatusPayload struct {
  Status    string `json:"status" binding:"required"`
}


var todoList = []Todo{}
// Note: Structのリテラルとしてmodel.Idだと使えないらしい
// var nextId model.Id
var nextId int

// LoadInitialData では規定のTodoの値を設定
func LoadInitialData() {
	todoList = []Todo{
		{ID: 1,	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: 2,	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: 3,	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: 4,	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: 5,	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}
	nextId = 6
}

// GetTodoList はGETでTODOリストを取得する
func GetTodoList(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, todoList)
	c.IndentedJSON(http.StatusOK, db.GetTodoList())
}

// GetTodoItemByID ではIDから任意のItemを取得する
func GetTodoItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	for _, item := range todoList {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}

// AddNewTodo では、POSTでItemを追加する
func AddNewTodo(c *gin.Context) {
	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: Payload"})
		return
	}

	newTodo := db.AddNewTodo(model.Payload{
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
	})
	/*
	newTodo := Todo {
		ID: nextId,
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
	}
	nextId +=  1

	todoList = append(todoList, newTodo)
	*/
	c.IndentedJSON(http.StatusCreated, newTodo)
}

// UpdateTodoItem ではIDで指定されたItemを更新する
func UpdateTodoItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: ID"})
		return
	}
	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: Payload"})
		return
	}

	newTodo := Todo {
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

// UpdateTodoState ではIDを指定したITEMのStatusを更新する
func UpdateTodoState(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: ID"})
		return
	}

	var payload StatusPayload
	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: Payload"})
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
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}

// DeleteTodoListItem ではIDで指定されたItemを削除する
func DeleteTodoListItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: ID"})
		return
	}
	for i, item := range todoList {
		if item.ID == id {
			todoList = append(todoList[:i], todoList[i + 1:]...)
			c.IndentedJSON(http.StatusOK, todoList)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo List Item not found"})
}
