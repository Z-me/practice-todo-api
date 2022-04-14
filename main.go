package main

import (
	todoApi "github.com/Z-me/practice-todo-api/api"
)

func main() {
	todoApi.Test()
	todoApi.Router().Run("localhost:8080")
	// todoApi.Router().Run()
}
