package main

import (
	"github.com/Z-me/practice-todo-api/api"
)

func main() {
	// db.ConnectDb()
	api.Test()
	api.Router().Run("localhost:8080")
}
