package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
)

func CreateUser(c *gin.Context) {
	var user models.UserPayload
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing password!"})
		return
	}
	user.Password = string(hashedPassword)

	db.Gorm.Create(&models.User{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	})

	c.Status(http.StatusCreated)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	db.Gorm.Find(&users)

	c.JSON(http.StatusOK, &users)
}

func GetUser(c *gin.Context) {
	var user models.User

	db.Gorm.Find(&user, c.Param("id"))

	c.JSON(http.StatusOK, &user)
}

func UpdateUser(c *gin.Context) {
	var user models.UserPayload
	var foundUser models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing password!"})
		return
	}
	user.Password = string(hashedPassword)

	db.Gorm.Find(&foundUser, c.Param("id"))

	foundUser.Name = user.Name
	foundUser.Email = user.Email
	foundUser.Password = user.Password

	db.Gorm.Save(&foundUser)

	c.JSON(http.StatusOK, &foundUser)
}

func DeleteUser(c *gin.Context) {
	db.Gorm.Delete(&models.User{}, c.Param("id"))

	c.Status(204)
}
