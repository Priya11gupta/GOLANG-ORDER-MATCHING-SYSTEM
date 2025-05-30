CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(20),
    side ENUM('buy', 'sell'),
    type ENUM('limit', 'market'),
    price DECIMAL(10,2),
    quantity INT,
    remaining_quantity INT,
    status ENUM('open', 'partial', 'filled', 'cancelled'),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trades (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(20),
    buy_order_id BIGINT,
    sell_order_id BIGINT,
    price DECIMAL(10,2),
    quantity INT,
    traded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
