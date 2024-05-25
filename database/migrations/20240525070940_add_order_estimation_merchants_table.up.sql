CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS order_estimation_merchants (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
order_estimation_id UUID NOT NULL REFERENCES order_estimation(id),
merchant_id UUID NOT NULL REFERENCES merchants(id),
is_starting_point BOOLEAN DEFAULT FALSE
);