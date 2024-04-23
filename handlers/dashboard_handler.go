package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
	cache "go-expense-tracker/redis"
)

func GetCurrentBalance(c *gin.Context) {
	DefaultDashboardFunc(c, "currentBalance", "")
	return
}

func GetTotalIncome(c *gin.Context) {
	DefaultDashboardFunc(c, "totalIncome", "value > 0")
	return
}

func GetTotalExpense(c *gin.Context) {
	DefaultDashboardFunc(c, "totalExpense", "value < 0")
	return
}

func DefaultDashboardFunc(c *gin.Context, cacheKey, query string) {
	var value int

	cachedVal, err := cache.Get(c, cacheKey)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"value": cachedVal})
		fmt.Println("em cache")
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

	GetTotalTransactionValue(query, value, userId.(uint))

	formattedValue := fmt.Sprintf("R$%.2f", float64(value)/100)

	err = cache.Set(c, "currentBalance", formattedValue, time.Minute*5)
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}

func GetTotalTransactionValue(query string, value int, userId uint) {
	accountsIds := db.Gorm.Model(models.Account{}).Where("user_id = ?", userId).Select("id")

	where := strings.Join([]string{"account_id in (?)", query}, " and ")

	db.Gorm.Model(&models.Transaction{}).
		Where(where, accountsIds).
		Select("sum(value)").
		Row().Scan(&value)
}
