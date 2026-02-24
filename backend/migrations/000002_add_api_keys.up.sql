CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    exchange VARCHAR(50) NOT NULL UNIQUE,
    api_key VARCHAR(500) NOT NULL,
    api_secret VARCHAR(500) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO api_keys (exchange, api_key, api_secret, is_active)
VALUES 
    ('mexc', '', '', false),
    ('bybit', '', '', false),
    ('gate', '', '', false),
    ('bitget', '', '', false)
ON CONFLICT (exchange) DO NOTHING;
