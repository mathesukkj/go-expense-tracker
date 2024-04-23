package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
	cache "go-expense-tracker/redis"
)

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	accountsIds := db.Gorm.Model(models.Account{}).Where("user_id = ?", userId).Select("id")

	db.Gorm.Where("account_id in (?)", accountsIds).Find(&transactions)

	c.JSON(http.StatusOK, &transactions)
}

func GetTransaction(c *gin.Context) {
	var transaction models.Transaction

	result := db.Gorm.First(&transaction, c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	err := checkUser(transaction.AccountID, userId.(uint))
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this account doesnt belong to you!"},
		)
		return
	}

	c.JSON(http.StatusOK, &transaction)
}

func CreateTransaction(c *gin.Context) {
	var transactionPayload models.TransactionPayload
	if err := c.ShouldBindJSON(&transactionPayload); err != nil {
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

	transaction := models.Transaction{
		Description:           transactionPayload.Description,
		Value:                 transactionPayload.Value,
		AccountID:             transactionPayload.AccountID,
		Date:                  transactionPayload.Date,
		Type:                  transactionPayload.Type,
		Completed:             transactionPayload.Completed,
		TransactionCategoryID: transactionPayload.TransactionCategoryID,
	}

	result := db.Gorm.Create(&transaction)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	err := checkUser(transaction.AccountID, userId.(uint))
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this account doesnt belong to you!"},
		)
		return
	}

	invalidateCache(c)

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func UpdateTransaction(c *gin.Context) {
	var transaction models.TransactionPayload
	var foundTransaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Gorm.Find(&foundTransaction, c.Param("id"))

	foundTransaction.Value = transaction.Value
	foundTransaction.Description = transaction.Description
	foundTransaction.Completed = transaction.Completed
	foundTransaction.Date = transaction.Date
	foundTransaction.Type = transaction.Type
	foundTransaction.TransactionCategoryID = transaction.TransactionCategoryID

	db.Gorm.Save(&foundTransaction)

	invalidateCache(c)

	c.JSON(http.StatusOK, &foundTransaction)
}

func DeleteTransaction(c *gin.Context) {
	db.Gorm.Delete(&models.Transaction{}, c.Param("id"))

	invalidateCache(c)

	c.Status(204)
}

func checkUser(accountId, userId uint) error {
	var account models.Account
	db.Gorm.Find(&account, accountId)

	if account.UserID != userId {
		return errors.New("account doesnt belonng to this user!")
	}
	return nil
}

func invalidateCache(c *gin.Context) {
	cache.Del(c, "currentBalance", "totalIncome", "totalExpense")
}
