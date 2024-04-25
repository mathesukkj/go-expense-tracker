package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/internal/db"
	"go-expense-tracker/internal/models"
)

func GetTransactionCategories(c *gin.Context) {
	var transactionCategories []models.TransactionCategory

	userId := c.MustGet("userId").(uint)

	db.Gorm.Find(&transactionCategories, "user_id = ?", userId)

	c.JSON(http.StatusOK, &transactionCategories)
}

func GetTransactionCategory(c *gin.Context) {
	var transactionCategory models.TransactionCategory

	if err := db.Gorm.First(&transactionCategory, c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction category not found"})
		return
	}

	userId := c.MustGet("userId").(uint)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, responseInvalidReqBody)
		return
	}

	userId := c.MustGet("userId").(uint)

	transactionCategory := models.TransactionCategory{
		UserID:   userId,
		Name:     transactionCategoryPayload.Name,
		ColorHex: transactionCategoryPayload.ColorHex,
		IconUrl:  transactionCategoryPayload.IconUrl,
	}

	if err := db.Gorm.Create(&transactionCategory); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func UpdateTransactionCategory(c *gin.Context) {
	var transactionCategory models.TransactionCategoryPayload
	var foundTransactionCategory models.TransactionCategory
	if err := c.ShouldBindJSON(&transactionCategory); err != nil {
		c.JSON(http.StatusBadRequest, responseInvalidReqBody)
		return
	}

	if err := db.Gorm.Find(&foundTransactionCategory, c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "transaction category not found"})
		return
	}

	userId := c.MustGet("userId").(uint)
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
	userId := c.MustGet("userId").(uint)

	var foundTransactionCategory models.TransactionCategory
	if err := db.Gorm.Find(&foundTransactionCategory, c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if userId != foundTransactionCategory.UserID {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "this category doesnt belong to you!"},
		)
		return
	}

	db.Gorm.Delete(foundTransactionCategory)

	c.Status(204)
}
