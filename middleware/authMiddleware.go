package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Z-me/practice-todo-api/lib/db"
	"github.com/Z-me/practice-todo-api/lib/util"
	"github.com/gin-gonic/gin"
)

func LoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Basic") {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			c.Abort()
			return
		}
		util.ConnectDB()
		defer util.DisconnectDB()
		dbObj := util.GetDbObj()

		fmt.Println("auth: ", auth)
		token := auth[6:]
		fmt.Println("token: ", token)
		splittedAuth := strings.Split(token, ":")
		name := splittedAuth[0]
		password := splittedAuth[1]
		if !db.CheckUserAuth(dbObj, name, password) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}
