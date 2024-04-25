package middleware

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you are not logged in!"})
			return
		}

		claims, err := getClaimsFromToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token!"})
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid user ID in token"},
			)
			return
		}

		c.Set("uid", uint(userID))
		c.Next()
	}
}

func getClaimsFromToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("SECRET_KEY"), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
