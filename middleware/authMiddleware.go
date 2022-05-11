package middleware

import (
	"net/http"
	"strings"

	"github.com/Z-me/practice-todo-api/lib/db"
	"github.com/gin-gonic/gin"
)

func LoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("auth")
		if auth == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			c.Abort()
			return
		}
		splittedAuth := strings.Split(auth, ":")
		name := splittedAuth[0]
		password := splittedAuth[1]
		if !db.CheckUserAuth(name, password) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}
