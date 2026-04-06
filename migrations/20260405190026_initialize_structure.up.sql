CREATE TABLE IF NOT EXISTS usdt_rates (
    id SERIAL PRIMARY KEY,
    ask_price NUMERIC(18, 2) NOT NULL,
    bid_price NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_usdt_rates_created_at ON usdt_rates(created_at);