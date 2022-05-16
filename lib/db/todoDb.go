package db

import (
	"time"

	"github.com/Z-me/practice-todo-api/api/model"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetNextID は次に指定するIDを取得する
func GetNextID(dbObj *gorm.DB) uint {
	todo := model.Todo{}
	dbObj.Last(&todo)
	return todo.ID + 1
}

// GetTodoList DBからTodoリストを取得して返却する関数
func GetTodoList(dbObj *gorm.DB) (model.TodoList, error) {
	todoList := model.TodoList{}
	err := dbObj.Find(&todoList).Error
	return todoList, err
}

// GetTodoItemByID はIDをもとにItemを取得する関数
func GetTodoItemByID(dbObj *gorm.DB, id uint) (model.Todo, error) {
	todo := model.Todo{}
	err := dbObj.First(&todo, id).Error
	return todo, err
}

// AddNewTodo はDBに指定のPayloadの値を投入
func AddNewTodo(dbObj *gorm.DB, payload model.Payload) (model.Todo, error) {
	newTodo := model.Todo{
		Title:     payload.Title,
		Status:    payload.Status,
		Details:   payload.Details,
		Priority:  payload.Priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result := dbObj.Create(&newTodo)

	return model.Todo{
		ID:        newTodo.ID,
		Title:     newTodo.Title,
		Status:    newTodo.Status,
		Details:   newTodo.Details,
		Priority:  newTodo.Priority,
		CreatedAt: newTodo.CreatedAt,
		UpdatedAt: newTodo.UpdatedAt,
	}, result.Error
}

// UpdateItem はDB上から指定のItemの情報を更新
func UpdateItem(dbObj *gorm.DB, id uint, payload model.Payload) (model.Todo, error) {

	target := model.Todo{}
	if err := dbObj.First(&target, id).Error; err != nil {
		return model.Todo{}, err
	}

	updated := model.Todo{
		ID:        id,
		Title:     payload.Title,
		Status:    payload.Status,
		Details:   payload.Details,
		Priority:  payload.Priority,
		CreatedAt: target.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := dbObj.Model(&target).Updates(&updated).Error; err != nil {
		return model.Todo{}, err
	}

	err := dbObj.First(&target, id).Error
	return updated, err
}

// UpdateItemStatus はDB上から指定のItemのStatusを更新
func UpdateItemStatus(dbObj *gorm.DB, id uint, status model.Status) (model.Todo, error) {
	target := model.Todo{}
	if err := dbObj.Find(&target, "id = ?", id).Error; err != nil {
		return model.Todo{}, err
	}

	if err := dbObj.Model(&target).Updates(map[string]interface{}{
		"Status":    status.Status,
		"UpdatedAt": time.Now(),
	}).Error; err != nil {
		return model.Todo{}, err
	}

	err := dbObj.First(&target, id).Error
	return target, err
}

// DeleteItem は任意のItemを削除
func DeleteItem(dbObj *gorm.DB, id uint) (model.Todo, error) {
	target := model.Todo{}
	if err := dbObj.Find(&target, "id = ?", id).Error; err != nil {
		return model.Todo{}, err
	}

	result := model.Todo{
		ID:       id,
		Title:    target.Title,
		Status:   target.Status,
		Details:  target.Details,
		Priority: target.Priority,
	}
	err := dbObj.Delete(&target).Error
	return result, err
}
