// api/handlers.go
package api

import (
	"net/http"

	"order-matching/db"
	"order-matching/models"
	"order-matching/services"

	"github.com/gin-gonic/gin"
)

func PlaceOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if order.Symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}
	if order.Side != "buy" && order.Side != "sell" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Side must be 'buy' or 'sell'"})
		return
	}
	if order.Type != "market" && order.Type != "limit" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'market' or 'limit'"})
		return
	}
	if order.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be positive"})
		return
	}
	if order.Type == "limit" && order.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be positive for limit orders"})
		return
	}

	// Insert order into database
	if err := models.InsertOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	// Try to match the order
	services.MatchOrder(&order)

	c.JSON(http.StatusCreated, order)
}

func CancelOrder(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("UPDATE orders SET status='cancelled' WHERE id=? AND status IN ('open', 'partial')", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to cancel"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

func GetOrderBook(c *gin.Context) {
	symbol := c.Query("symbol")
	book := gin.H{"bids": []models.Order{}, "asks": []models.Order{}}

	buyRows, _ := db.DB.Query("SELECT price, remaining_quantity FROM orders WHERE symbol=? AND side='buy' AND status IN ('open', 'partial') ORDER BY price DESC", symbol)
	for buyRows.Next() {
		var price float64
		var qty int
		buyRows.Scan(&price, &qty)
		book["bids"] = append(book["bids"].([]models.Order), models.Order{Price: price, RemainingQty: qty})
	}

	sellRows, _ := db.DB.Query("SELECT price, remaining_quantity FROM orders WHERE symbol=? AND side='sell' AND status IN ('open', 'partial') ORDER BY price ASC", symbol)
	for sellRows.Next() {
		var price float64
		var qty int
		sellRows.Scan(&price, &qty)
		book["asks"] = append(book["asks"].([]models.Order), models.Order{Price: price, RemainingQty: qty})
	}

	c.JSON(http.StatusOK, book)
}

func GetTrades(c *gin.Context) {
	symbol := c.Query("symbol")
	rows, _ := db.DB.Query(`
        SELECT id, symbol, buy_order_id, sell_order_id, price, quantity, traded_at 
        FROM trades 
        WHERE symbol=? 
        ORDER BY traded_at DESC`,
		symbol,
	)
	var trades []models.Trade
	for rows.Next() {
		var t models.Trade
		rows.Scan(
			&t.ID,
			&t.Symbol,
			&t.BuyOrderID,
			&t.SellOrderID,
			&t.Price,
			&t.Quantity,
			&t.TradedAt,
		)
		trades = append(trades, t)
	}
	c.JSON(http.StatusOK, trades)
}
