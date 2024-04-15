package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
)

func GetAccounts(c *gin.Context) {
	var accounts []models.Account

	db.Gorm.Find(&accounts)

	c.JSON(http.StatusOK, &accounts)
}

func GetAccount(c *gin.Context) {
	var account models.Account

	result := db.Gorm.Find(&account, c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "account not found"})
	}
	c.JSON(http.StatusOK, &account)
}

func CreateAccount(c *gin.Context) {
	var accountPayload models.AccountPayload
	if err := c.ShouldBindJSON(&accountPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	account := models.Account{
		UserID:  userId.(uint),
		Name:    accountPayload.Name,
		Balance: accountPayload.Balance,
	}

	result := db.Gorm.Create(&account)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func UpdateAccount(c *gin.Context) {
	var account models.AccountPayload
	var foundAccount models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Gorm.Find(&foundAccount, c.Param("id"))

	foundAccount.Name = account.Name
	foundAccount.Balance = account.Balance

	db.Gorm.Save(&foundAccount)

	c.JSON(http.StatusOK, &foundAccount)
}

func DeleteAccount(c *gin.Context) {
	db.Gorm.Delete(&models.Account{}, c.Param("id"))

	c.Status(204)
}
