package main

import (
	"github.com/Z-me/practice-todo-api/middleware"
)

func main() {
	middleware.ConnectDb()
	// api.Test()
	// api.Router().Run("localhost:8080")
	// todoApi.Router().Run()
}
