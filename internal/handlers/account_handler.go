package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/internal/db"
	"go-expense-tracker/internal/models"
)

func GetAccounts(c *gin.Context) {
	var accounts []models.Account

	userId := c.MustGet("userId").(uint)

	db.Gorm.Find(&accounts, "user_id = ?", userId)

	c.JSON(http.StatusOK, &accounts)
}

func GetAccount(c *gin.Context) {
	var account models.Account

	userId := c.MustGet("userId").(uint)

	result := db.Gorm.First(&account, c.Param("id"), "user_id = ?", userId)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse("account"))
		return
	}

	c.JSON(http.StatusOK, &account)
}

func CreateAccount(c *gin.Context) {
	var accountPayload models.AccountPayload
	if err := c.ShouldBindJSON(&accountPayload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responseInvalidReqBody)
		return
	}

	userId := c.MustGet("userId").(uint)

	account := models.Account{
		UserID:  userId,
		Name:    accountPayload.Name,
		Balance: accountPayload.Balance,
	}

	if err := db.Gorm.Create(&account).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func UpdateAccount(c *gin.Context) {
	var account models.AccountPayload
	var foundAccount models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, responseInvalidReqBody)
		return
	}

	userId := c.MustGet("userId").(uint)

	if err := db.Gorm.Find(&foundAccount, c.Param("id"), "user_id = ?", userId).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse("account"))
		return
	}

	foundAccount.Name = account.Name
	foundAccount.Balance = account.Balance

	db.Gorm.Save(&foundAccount)

	c.JSON(http.StatusOK, &foundAccount)
}

func DeleteAccount(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	if err := db.Gorm.Delete(&models.Account{}, c.Param("id"), "user_id = ?", userId).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse("account"))
		return
	}

	c.Status(204)
}
