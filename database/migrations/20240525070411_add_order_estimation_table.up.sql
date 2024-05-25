CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS order_estimation (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id),
user_location_lat REAL NOT NULL,
user_location_long REAL NOT NULL,
total_price INTEGER NOT NULL,
estimated_delivery_time INTEGER NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);