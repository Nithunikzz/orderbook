package handlers

import (
	"net/http"
	"orderbook/models"
	"strconv"

	"orderbook/orderbook"

	"github.com/gin-gonic/gin"
)

func PlaceOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := models.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	go orderbook.MatchOrders(order.Pair)

	c.JSON(http.StatusOK, order)
}

func GetOrderBook(c *gin.Context) {
	pair := c.Query("pair")
	depthParam := c.DefaultQuery("depth", "10")

	depth, err := strconv.Atoi(depthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid depth"})
		return
	}

	buyOrders := models.GetOpenOrders(pair, "buy")
	sellOrders := models.GetOpenOrders(pair, "sell")

	// Sorting in price/time priority
	buySorted := sortOrders(buyOrders, true, depth)
	sellSorted := sortOrders(sellOrders, false, depth)

	c.JSON(http.StatusOK, gin.H{
		"buy":  buySorted,
		"sell": sellSorted,
	})
}

func sortOrders(orders []models.Order, desc bool, depth int) []gin.H {
	sorted := make([]gin.H, 0)
	count := 0

	for _, o := range orders {
		sorted = append(sorted, gin.H{
			"price":    o.Price,
			"quantity": o.Quantity - o.FilledQuantity,
		})
		count++
		if count >= depth {
			break
		}
	}

	return sorted
}

func GetUserOrders(c *gin.Context) {
	uidStr := c.Query("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	orders := models.GetUserOrders(uid)
	c.JSON(http.StatusOK, orders)
}
func GetTotalFees(c *gin.Context) {
	var total float64
	err := models.DB.QueryRow("SELECT SUM(fee) FROM orders").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total fee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total fee is collected": total})
}
