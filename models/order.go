package models

import (
	"fmt"
	"order-matching/db"
	"time"
)

type Order struct {
	ID           int       `json:"id"`
	Symbol       string    `json:"symbol"`
	Side         string    `json:"side"`
	Type         string    `json:"type"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	RemainingQty int       `json:"remaining_quantity"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type Trade struct {
	ID          int       `json:"id"`
	Symbol      string    `json:"symbol"`
	BuyOrderID  int       `json:"buy_order_id"`
	SellOrderID int       `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	TradedAt    time.Time `json:"traded_at"`
}

func InsertOrder(order *Order) error {
	order.Status = "open"
	order.RemainingQty = order.Quantity
	order.CreatedAt = time.Now()

	query := `
        INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := db.DB.Exec(query,
		order.Symbol, order.Side, order.Type,
		order.Price, order.Quantity, order.RemainingQty,
		order.Status, order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}
	order.ID = int(id)

	return nil
}

func UpdateOrder(order *Order) error {
	query := `UPDATE orders SET remaining_quantity=?, status=? WHERE id=?`
	_, err := db.DB.Exec(query, order.RemainingQty, order.Status, order.ID)
	return err
}

func InsertTrade(trade *Trade) error {
	query := `INSERT INTO trades (symbol, buy_order_id, sell_order_id, price, quantity, traded_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.DB.Exec(query, trade.Symbol, trade.BuyOrderID, trade.SellOrderID, trade.Price, trade.Quantity, trade.TradedAt)
	return err
}
