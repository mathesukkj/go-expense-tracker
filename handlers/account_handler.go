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

	db.Gorm.Find(&account, c.Param("id"))

	c.JSON(http.StatusOK, &account)
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
