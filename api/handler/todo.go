package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Z-me/practice-todo-api/api/model"
	db "github.com/Z-me/practice-todo-api/middleware"
)

// Todo APIのレスポンスの構造体
type Todo struct {
  ID        int     `json:"id"`
  Title     string  `json:"title" binding:"required,max=30"`
  Status    string  `json:"status" binding:"required"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority" binding:"required,max=1000"`
}

// Payload APIのDBの新規作成及び更新のPayload
type Payload struct {
  Title     string  `json:"title" binding:"required,max=30"`
  Status    string  `json:"status" binding:"required"`
  Details   string  `json:"details"`
  Priority  string  `json:"priority" binding:"required,max=1000"`
}

// StatusPayload APIのStatusのみ更新する際のPayload
type StatusPayload struct {
  Status    string `json:"status" binding:"required"`
}

// GetTodoList はGETでTODOリストを取得する
func GetTodoList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.GetTodoList(c))
}

// GetTodoItemByID ではIDから任意のItemを取得する
func GetTodoItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: id"})
	}

	c.IndentedJSON(http.StatusOK, db.GetTodoItemByID(c, uint(id)))
}

// AddNewTodo では、POSTでItemを追加する
func AddNewTodo(c *gin.Context) {
	var payload Payload

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: Payload"})
		return
	}

	newTodo := db.AddNewTodo(c,
		model.Payload{
			Title: payload.Title,
			Status: payload.Status,
			Details: payload.Details,
			Priority: payload.Priority,
		})
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

	updated := db.UpdateItem(c, uint(id), model.Payload{
		Title: payload.Title,
		Status: payload.Status,
		Details: payload.Details,
		Priority: payload.Priority,
	})
	c.IndentedJSON(http.StatusCreated, updated)
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

	updated := db.UpdateItemStatus(c, uint(id), model.Status{
		Status: payload.Status,
	})
	c.IndentedJSON(http.StatusCreated, updated)
}

// DeleteTodoListItem ではIDで指定されたItemを削除する
func DeleteTodoListItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request: ID"})
		return
	}

	deleted := db.DeleteItem(c, uint(id))
	c.IndentedJSON(http.StatusOK, deleted)
}
