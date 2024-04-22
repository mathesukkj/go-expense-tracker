package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
)

func GetCurrentBalance(c *gin.Context) {
	var currentBalance int

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	GetTotalTransactionValue("", currentBalance, userId.(uint))

	c.JSON(http.StatusOK, gin.H{"value": currentBalance})
}

func GetTotalIncome(c *gin.Context) {
	var totalIncome int

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	GetTotalTransactionValue("value > 0", totalIncome, userId.(uint))

	c.JSON(http.StatusOK, gin.H{"value": totalIncome})
}

func GetTotalExpense(c *gin.Context) {
	var totalExpense int

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	GetTotalTransactionValue("value < 0", totalExpense, userId.(uint))

	c.JSON(http.StatusOK, gin.H{"value": totalExpense})
}

func GetTotalTransactionValue(query string, value int, userId uint) {
	accountsIds := db.Gorm.Model(models.Account{}).Where("user_id = ?", userId).Select("id")

	where := strings.Join([]string{"account_id in (?)", query}, " and ")

	db.Gorm.Model(&models.Transaction{}).
		Where(where, accountsIds).
		Select("sum(value)").
		Row().Scan(&value)
}
