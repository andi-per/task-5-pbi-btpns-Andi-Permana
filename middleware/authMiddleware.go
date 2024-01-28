package middleware

import (
	"api/rakamin-api/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateToken(c *gin.Context) {
	tokenString := helpers.ExtractTokenFromHeader(c.Request)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
		c.Abort()
		return
	}

	// Parse and validate the JWT token
	claims, err := helpers.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Set the user ID and email from the claims in the context for later use
	c.Set("userID", claims["user_id"])
	c.Set("email", claims["email"])

	// Continue processing the request
	c.Next()
}