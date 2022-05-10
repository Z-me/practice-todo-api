package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/model"
	db "github.com/Z-me/practice-todo-api/lib"
)

// Todo APIのレスポンスの構造体
type Todo struct {
	ID			int     	`json:"id"`
	Title		string  	`json:"title" binding:"required,max=30"`
	Status		string  	`json:"status" binding:"required"`
	Details		string  	`json:"details"`
	Priority	string  	`json:"priority" binding:"required,max=1000"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
}

// Payload APIのDBの新規作成及び更新のPayload
type Payload struct {
  	Title     	string  `json:"title" binding:"required,max=30"`
  	Status    	string  `json:"status" binding:"required"`
  	Details   	string  `json:"details"`
  	Priority  	string  `json:"priority" binding:"required,max=1000"`
}

// StatusPayload APIのStatusのみ更新する際のPayload
type StatusPayload struct {
  	Status		string `json:"status" binding:"required"`
}

func connectDB(c *gin.Context){
	if err := db.ConnectDB(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to connect database"})
	}
}

// GetTodoList はGETでTODOリストを取得する
func GetTodoList(c *gin.Context) {
	connectDB(c)
	defer db.DisconnectDB()

	todoList, err := db.GetTodoList();
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Todo List Item not found"})
	}
	result := []Todo{}
	for _, v := range todoList{
		item := Todo{
			ID:			int(v.ID),
			Title:		v.Title,
			Status:		v.Status,
			Details:	v.Details,
			Priority:	v.Priority,
			CreatedAt:	v.CreatedAt,
			UpdatedAt:	v.UpdatedAt,
		}
		result = append(result, item)
	}
	c.IndentedJSON(http.StatusOK, result)
}

// GetTodoItemByID ではIDから任意のItemを取得する
func GetTodoItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: id"})
	}

	connectDB(c)
	defer db.DisconnectDB()

	item, err := db.GetTodoItemByID(uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Target item is not found"})
	}
	c.IndentedJSON(http.StatusOK, Todo{
		ID:			int(item.ID),
		Title:		item.Title,
		Status:		item.Status,
		Details:	item.Details,
		Priority:	item.Priority,
		CreatedAt:	item.CreatedAt,
		UpdatedAt:	item.UpdatedAt,
	})
}

// AddNewTodo では、POSTでItemを追加する
func AddNewTodo(c *gin.Context) {
	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: Payload"})
		return
	}

	connectDB(c)
	defer db.DisconnectDB()

	newTodo, err := db.AddNewTodo(
		model.Payload{
			Title: payload.Title,
			Status: payload.Status,
			Details: payload.Details,
			Priority: payload.Priority,
		})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to create new item"})
	}
	c.IndentedJSON(http.StatusCreated, Todo{
		ID:			int(newTodo.ID),
		Title:		newTodo.Title,
		Status:		newTodo.Status,
		Details:	newTodo.Details,
		Priority:	newTodo.Priority,
		CreatedAt:	newTodo.CreatedAt,
		UpdatedAt:	newTodo.UpdatedAt,
	})
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

	connectDB(c)
	defer db.DisconnectDB()

	updated, err := db.UpdateItem(uint(id), model.Payload{
		Title:		payload.Title,
		Status:		payload.Status,
		Details:	payload.Details,
		Priority:	payload.Priority,
	})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to update item"})
	}
	fmt.Println("updated", updated)
	c.IndentedJSON(http.StatusOK, Todo{
		ID:			id,
		Title:		updated.Title,
		Status:		updated.Status,
		Details:	updated.Details,
		Priority:	updated.Priority,
		CreatedAt:	updated.CreatedAt,
		UpdatedAt:	updated.UpdatedAt,
	})
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

	connectDB(c)
	defer db.DisconnectDB()

	updated, err := db.UpdateItemStatus(uint(id), model.Status{
		Status: payload.Status,
	})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to update item"})
	}
	c.IndentedJSON(http.StatusOK, updated)
}

// DeleteTodoListItem ではIDで指定されたItemを削除する
func DeleteTodoListItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: ID"})
		return
	}

	connectDB(c)
	defer db.DisconnectDB()

	deleted, err := db.DeleteItem(uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "fail to update item"})
	}

	c.IndentedJSON(http.StatusOK, deleted)
}
