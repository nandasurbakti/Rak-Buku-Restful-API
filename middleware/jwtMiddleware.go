package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nandasurbakti/Rak-Buku-Restful-API/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}
	
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error(),})
		c.Abort()
		return 
	}

	// userID := claims.Username
	// if userID == ""{
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan"})
	// 	c.Abort()
	// 	return
	// }

	// c.Set("userID", userID)
	c.Set("username", claims.Username)
	c.Next()
	}

}