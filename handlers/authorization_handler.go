package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
)

var secretKey = []byte("my-secret-string")

func Signup(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var signupData models.UserPayload
	if err := c.ShouldBindJSON(&signupData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(signupData.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing password!"})
		return
	}
	signupData.Password = string(hashedPassword)

	newUser := models.User{
		Name:     signupData.Name,
		Password: signupData.Password,
		Email:    signupData.Email,
	}

	db.Gorm.Create(&newUser)

	token, err := createToken(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "thats weird. an error happened!"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged in successfully"})
}

func Login(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var loginData models.LoginPayload
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db.Gorm.Find(&user, "email = ?", loginData.Email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong email or password"})
		return
	}

	token, err := createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "thats weird. an error happened!"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged in successfully"})
}

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
		return secretKey, nil
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

func createToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	return token.SignedString(secretKey)
}
