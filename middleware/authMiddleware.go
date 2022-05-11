package middleware

import (
	"fmt"
	"net/http"
	"strings"

	db "github.com/Z-me/practice-todo-api/lib/db"
	"github.com/gin-gonic/gin"
)

func LoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("aaaaaaaaaa")
		auth := c.Request.Header.Get("auth")
		// name := c.Request.Header.Get("name")
		// password := c.Request.Header.Get("password")
		splittedAuth := strings.Split(auth, ":")
		name := splittedAuth[0]
		password := splittedAuth[1]
		if name == "" || password == "" || !db.CheckUserAuth(name, password) {
			// c.Status(http.StatusUnauthorized)
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}
