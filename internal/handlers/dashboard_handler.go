package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/internal/db"
	"go-expense-tracker/internal/models"
	cache "go-expense-tracker/internal/redis"
	datefmt "go-expense-tracker/internal/utils"
)

func GetCurrentBalance(c *gin.Context) {
	DefaultDashboardFunc(c, "currentBalance", "")
	return
}

func GetMonthlyBalance(c *gin.Context) {
	firstDay := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)

	query := "date between " + datefmt.FormatDateToPostgres(
		firstDay,
	) + " and " + datefmt.FormatDateToPostgres(
		lastDay,
	)

	DefaultDashboardFunc(c, "monthlyBalance", query)
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
		fmt.Println(cacheKey, "em cache")
		c.JSON(http.StatusOK, gin.H{"value": cachedVal})
		return
	}

	userId := c.MustGet("userId").(uint)

	GetTotalTransactionValue(query, value, userId)

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
