package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
)

func Login(c *gin.Context) {
	var loginData models.LoginPayload
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db.Gorm.Find(&user, "email = ?", loginData.Email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "passwords are not equal"})
		return
	}

	c.String(200, "fala paizaao")
}
