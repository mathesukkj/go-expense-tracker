package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/internal/db"
	"go-expense-tracker/internal/models"
	cache "go-expense-tracker/internal/redis"
)

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	userId := c.MustGet("userId").(uint)

	dbQuery := db.Gorm.Model(&models.Transaction{})

	if category := c.Query("category_id"); category != "" {
		id, err := strconv.Atoi(category)
		if err != nil {
			c.JSON(400, gin.H{"error": "category_id must be a number"})
			return
		}
		dbQuery = dbQuery.Where("transaction_category_id = ?", id)
	}

	if account := c.Query("account_id"); account != "" {
		id, err := strconv.Atoi(account)
		if err != nil {
			c.JSON(400, gin.H{"error": "account_id must be a number"})
			return
		}
		dbQuery = dbQuery.Where("account_id = ?", id)
	} else {
		accountsIds := db.Gorm.Model(models.Account{}).Where("user_id = ?", userId).Select("id")
		dbQuery = dbQuery.Where("account_id in (?)", accountsIds)
	}

	if transactionType := c.Query("transaction_type"); transactionType != "" {
		switch transactionType {
		case "expense":
			dbQuery = dbQuery.Where("value < 0")
		case "income":
			dbQuery = dbQuery.Where("value > 0")
		}
	}

	if search := c.Query("search"); search != "" {
		dbQuery = dbQuery.Where("description LIKE ?", "%"+search+"%")
	}

	dbQuery.Find(&transactions)

	c.JSON(http.StatusOK, &transactions)
}

func GetTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := db.Gorm.First(&transaction, c.Param("id")).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	userId := c.MustGet("userId").(uint)

	if err := checkUser(transaction.AccountID, userId); err != nil {
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
		c.AbortWithStatusJSON(http.StatusBadRequest, responseInvalidReqBody)
		return
	}

	userId := c.MustGet("userId").(uint)

	transaction := models.Transaction{
		Description:           transactionPayload.Description,
		Value:                 transactionPayload.Value,
		AccountID:             transactionPayload.AccountID,
		Date:                  transactionPayload.Date,
		Type:                  transactionPayload.Type,
		Completed:             transactionPayload.Completed,
		TransactionCategoryID: transactionPayload.TransactionCategoryID,
	}

	if err := db.Gorm.Create(&transaction).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to create transaction"})
		return
	}

	if err := checkUser(transaction.AccountID, userId); err != nil {
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
		c.JSON(http.StatusBadRequest, responseInvalidReqBody)
		return
	}

	if err := db.Gorm.Find(&foundTransaction, c.Param("id")).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, notFoundResponse("transaction"))
		return
	}

	userId := c.MustGet("userId").(uint)
	if err := checkUser(transaction.AccountID, userId); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this account doesnt belong to you!"},
		)
		return
	}

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
	userId := c.MustGet("userId").(uint)

	var transaction models.Transaction
	if err := db.Gorm.First(&transaction, c.Param("id")).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	if err := checkUser(transaction.AccountID, userId); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this account doesnt belong to you!"},
		)
		return
	}

	db.Gorm.Delete(transaction)

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
