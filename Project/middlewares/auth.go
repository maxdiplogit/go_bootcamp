package middlewares

import (
	"net/http"

	"example.com/app/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})
		return
	}

	userID, err := utils.VerifyJWTToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized, JWT token invalid"})
		return
	}

	context.Set("userID", userID)

	context.Next()
}
