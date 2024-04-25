package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-expense-tracker/internal/db"
	"go-expense-tracker/internal/models"
)

func Signup(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.JSON(http.StatusUnauthorized, responseAlreadyLoggedIn)
		return
	}

	var signupData models.UserPayload
	if err := c.ShouldBindJSON(&signupData); err != nil {
		c.JSON(http.StatusBadRequest, responseInvalidReqBody)
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
		c.JSON(http.StatusUnauthorized, responseAlreadyLoggedIn)
		return
	}

	var loginData models.LoginPayload
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, responseInvalidReqBody)
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

func Signout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "signed out succesfully"})
}

func createToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	return token.SignedString(os.Getenv("SECRET_KEY"))
}
