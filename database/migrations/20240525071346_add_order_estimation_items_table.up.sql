CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS order_estimation_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_estimation_id UUID NOT NULL REFERENCES order_estimation(id),
    item_id UUID NOT NULL REFERENCES items(id),
    quantity INTEGER NOT NULL
);