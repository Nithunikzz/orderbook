package main

import (
	"log"

	"orderbook/handlers"
	"orderbook/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDB()

	r := gin.Default()
	r.POST("/api/pairs", handlers.AddCurrencyPair)
	r.POST("/api/orders", handlers.PlaceOrder)
	r.GET("/api/orderbook", handlers.GetOrderBook)
	r.GET("/api/orders", handlers.GetUserOrders)
	r.GET("/api/fees/total", handlers.GetTotalFees)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}

// curl -X POST http://localhost:8080/api/pairs \
//   -H "Content-Type: application/json" \
//   -d '{"pair": "BTC/USDT"}'
//Invoke-WebRequest -Uri http://localhost:8080/api/pairs -Method POST -Body '{"base":"EXT", "quote":"USDT"}' -Headers @{ "Content-Type" = "application/json" }
//Invoke-WebRequest -Uri http://localhost:8080/api/orders -Method POST -Body '{"user_id":1,"pair":"ETX/USDT","price":28000,"quantity":0.5,"side":"buy"}' -Headers @{ "Content-Type" = "application/json" }
//Invoke-WebRequest -Uri http://localhost:8080/api/orders -Method POST -Body '{"user_id":2,"pair":"ETX/USDT","price":26000,"quantity":0.4,"side":"sell"}' -Headers @{ "Content-Type" = "application/json" }
//Invoke-WebRequest -Uri "http://localhost:8080/api/orderbook?pair=BTC/USDT&limit=5"
//Invoke-WebRequest -Uri "http://localhost:8080/api/orders?user_id=1"
