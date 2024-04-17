package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
	"go-expense-tracker/models"
)

func GetTransactionCategories(c *gin.Context) {
	var transactionCategories []models.TransactionCategory

	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	db.Gorm.Find(&transactionCategories, "user_id = ?", userId)

	c.JSON(http.StatusOK, &transactionCategories)
}

func GetTransactionCategory(c *gin.Context) {
	var transactionCategory models.TransactionCategory

	result := db.Gorm.First(&transactionCategory, c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction category not found"})
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

	if userId != transactionCategory.UserID {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this category doesnt belong to you!"},
		)
		return
	}

	c.JSON(http.StatusOK, &transactionCategory)
}

func CreateTransactionCategory(c *gin.Context) {
	var transactionCategoryPayload models.TransactionCategoryPayload
	if err := c.ShouldBindJSON(&transactionCategoryPayload); err != nil {
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

	transactionCategory := models.TransactionCategory{
		UserID:   userId.(uint),
		Name:     transactionCategoryPayload.Name,
		ColorHex: transactionCategoryPayload.ColorHex,
		IconUrl:  transactionCategoryPayload.IconUrl,
	}

	result := db.Gorm.Create(&transactionCategory)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	if transactionCategory.UserID != userId {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this account doesnt belong to you!"},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func UpdateTransactionCategory(c *gin.Context) {
	var transactionCategory models.TransactionCategoryPayload
	var foundTransactionCategory models.TransactionCategory
	if err := c.ShouldBindJSON(&transactionCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Gorm.Find(&foundTransactionCategory, c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction category not found"})
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

	if userId != foundTransactionCategory.UserID {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this category doesnt belong to you!"},
		)
		return
	}

	foundTransactionCategory.Name = transactionCategory.Name
	foundTransactionCategory.ColorHex = transactionCategory.ColorHex
	foundTransactionCategory.IconUrl = transactionCategory.IconUrl

	db.Gorm.Save(&foundTransactionCategory)

	c.JSON(http.StatusOK, &foundTransactionCategory)
}

func DeleteTransactionCategory(c *gin.Context) {
	userId, ok := c.Get("uid")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found. please login again"},
		)
		return
	}

	var foundTransactionCategory models.TransactionCategory
	result := db.Gorm.Find(&foundTransactionCategory, c.Param("id"))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction category not found"})
		return
	}

	if userId != foundTransactionCategory.UserID {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this category doesnt belong to you!"},
		)
		return
	}

	db.Gorm.Delete(&models.TransactionCategory{}, c.Param("id"))

	c.Status(204)
}
