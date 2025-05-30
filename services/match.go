package services

import (
	"log"
	"order-matching/db"
	"order-matching/models"
	"time"
)

func MatchOrder(newOrder *models.Order) {
	// Get matching orders from opposite side
	opposite := "buy"
	priceCond := ">="
	orderBy := "price DESC"
	if newOrder.Side == "buy" {
		opposite = "sell"
		priceCond = "<="
		orderBy = "price ASC"
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer tx.Rollback()

	// Query matching orders
	query := `SELECT id, price, remaining_quantity FROM orders 
		WHERE symbol = ? AND side = ? AND price ` + priceCond + ` ? 
		AND status IN ('open', 'partial') 
		ORDER BY ` + orderBy + `, created_at ASC`

	rows, err := tx.Query(query, newOrder.Symbol, opposite, newOrder.Price)
	if err != nil {
		log.Printf("Query error: %v", err)
		return
	}
	defer rows.Close()

	remaining := newOrder.RemainingQty
	matched := false

	// Match with existing orders
	for rows.Next() && remaining > 0 {
		var id int
		var price float64
		var availQty int
		err = rows.Scan(&id, &price, &availQty)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		// Calculate trade quantity
		matchQty := min(remaining, availQty)
		remaining -= matchQty
		matched = true

		// Record trade
		trade := &models.Trade{
			Symbol:   newOrder.Symbol,
			Price:    price,
			Quantity: matchQty,
			TradedAt: time.Now(),
		}
		if newOrder.Side == "buy" {
			trade.BuyOrderID = newOrder.ID
			trade.SellOrderID = id
		} else {
			trade.SellOrderID = newOrder.ID
			trade.BuyOrderID = id
		}

		// Insert trade record
		_, err = tx.Exec(
			"INSERT INTO trades (symbol, buy_order_id, sell_order_id, price, quantity, traded_at) VALUES (?, ?, ?, ?, ?, ?)",
			trade.Symbol, trade.BuyOrderID, trade.SellOrderID, trade.Price, trade.Quantity, trade.TradedAt,
		)
		if err != nil {
			log.Printf("Failed to insert trade: %v", err)
			return
		}

		// Update matched order
		newAvail := availQty - matchQty
		status := "partial"
		if newAvail == 0 {
			status = "filled"
		}
		_, err = tx.Exec(
			"UPDATE orders SET remaining_quantity = ?, status = ? WHERE id = ?",
			newAvail, status, id,
		)
		if err != nil {
			log.Printf("Failed to update matched order: %v", err)
			return
		}
	}

	// Update new order status
	newOrder.RemainingQty = remaining
	if remaining == 0 {
		newOrder.Status = "filled"
	} else {
		if newOrder.Type == "market" {
			// Market orders should be cancelled if not fully filled
			newOrder.Status = "cancelled"
			newOrder.RemainingQty = 0
		} else if matched {
			// Limit orders should be partial if partially filled
			newOrder.Status = "partial"
		} else {
			// Limit orders should be open if not matched
			newOrder.Status = "open"
		}
	}

	_, err = tx.Exec(
		"UPDATE orders SET remaining_quantity = ?, status = ? WHERE id = ?",
		newOrder.RemainingQty, newOrder.Status, newOrder.ID,
	)
	if err != nil {
		log.Printf("Failed to update new order: %v", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
