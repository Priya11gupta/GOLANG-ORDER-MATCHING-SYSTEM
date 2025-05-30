package main

import (
	"log"

	"order-matching/api"
	"order-matching/db"

	"github.com/gin-gonic/gin"
)

func main() {
	//Connect to MySQL
	db.ConnectDB()

	//Initialize Gin router
	router := gin.Default()

	//Define API routes
	router.POST("/orders", api.PlaceOrder)        // Place a new order
	router.DELETE("/orders/:id", api.CancelOrder) // Cancel an order by ID
	router.GET("/orderbook", api.GetOrderBook)    // Get current order book for a symbol
	router.GET("/trades", api.GetTrades)          // Get all trades for a symbol
	//router.GET("/orders/:id", api.GetOrder)       // Get order details by ID

	//Start the server
	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
