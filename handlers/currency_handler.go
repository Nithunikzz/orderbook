package handlers

import (
	"net/http"
	"orderbook/models"

	"github.com/gin-gonic/gin"
)

func AddCurrencyPair(c *gin.Context) {
	var pair models.CurrencyPair
	if err := c.ShouldBindJSON(&pair); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := models.AddPair(pair)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add pair"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Currency pair added"})
}
