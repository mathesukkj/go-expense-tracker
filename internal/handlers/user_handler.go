package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-expense-tracker/internal/db"
	"go-expense-tracker/internal/models"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	db.Gorm.Find(&users)

	c.JSON(http.StatusOK, &users)
}

func GetUser(c *gin.Context) {
	var user models.User

	result := db.Gorm.First(&user, c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
	}

	c.JSON(http.StatusOK, &user)
}

func UpdateUser(c *gin.Context) {
	var user models.UserPayload
	var foundUser models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responseInvalidReqBody)
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
