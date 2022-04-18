package main

import (
	"github.com/Z-me/practice-todo-api/api"
)

func main() {
	api.Test()
	api.Router().Run("localhost:8080")
	// todoApi.Router().Run()
}
