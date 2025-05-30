# Order Matching Engine

A simplified order matching system implemented in Go that matches buy and sell orders based on price-time priority. The system supports both market and limit orders, with a RESTful API interface and MySQL persistence.

## Features

- Support for both Limit and Market orders
- Price-time priority matching algorithm
- RESTful API for order management
- MySQL persistence for orders and trades
- Real-time order book maintenance
- Support for partial fills
- Proper error handling and validation

## Prerequisites

- Go 1.24 or later
- MySQL 8.0 or later
- Git

## Setup

1. Clone the repository:
```bash
git clone https://github.com/Priya11gupta/GOLANG-ORDER-MATCHING-SYSTEM.git
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the MySQL database:
```bash
mysql -u root -p
```

```sql
CREATE DATABASE order_matching;
USE order_matching;
```

4. Run the database migrations:
```bash
mysql -u root -p order_matching < migrations/schema.sql
```

5. Configure the database connection in `db/mysql.go` if needed.

6. Build and run the application:
```bash
go run main.go
```

The server will start on http://localhost:8080

## API Endpoints

### Place Order
```
POST /orders
```
Request body:
```json
{
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "price": 150.00,
    "quantity": 100
}
```

### Cancel Order
```
DELETE /orders/{orderId}
```

### Get Order Book
```
GET /orderbook?symbol=AAPL
```

### List Trades
```
GET /trades?symbol=AAPL
```

### Get Order Status
```
GET /orders/{orderId}
```

## Example Usage

1. Place a limit buy order:
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "price": 150.00,
    "quantity": 100
  }'
```

2. Place a market sell order:
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "side": "sell",
    "type": "market",
    "quantity": 50
  }'
```

3. Get the order book:
```bash
curl http://localhost:8080/orderbook?symbol=AAPL
```

4. Cancel an order:
```bash
curl -X DELETE http://localhost:8080/orders/1
```

## Design Decisions

1. In-Memory Order Book
   - Maintains a sorted list of orders for quick matching
   - Synchronized with database for persistence
   - Uses read-write mutex for thread safety

2. Database Schema
   - Orders table tracks all order states
   - Trades table records matched executions
   - No ORM used, raw SQL for better performance

3. Price-Time Priority
   - Orders at better prices are matched first
   - For same price, older orders get priority (FIFO)

4. Market Orders
   - Execute immediately at best available price
   - Partial fills possible
   - Remaining quantity is canceled if no matches

5. Error Handling
   - Proper validation of all inputs
   - Appropriate HTTP status codes
   - Detailed error messages
   - Transaction support for data consistency

## Limitations and Potential Improvements

1. Performance
   - Could implement price level aggregation
   - Could use memory-mapped files for persistence
   - Could add caching layer

2. Features
   - Add support for stop orders
   - Add support for IOC/FOK orders
   - Add WebSocket for real-time updates

3. Scalability
   - Could implement horizontal scaling
   - Could add message queue for order processing
   - Could separate matching engine from API layer
   