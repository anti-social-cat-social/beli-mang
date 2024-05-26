CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE product_categories
 AS ENUM (
'Beverage',
'Food',
'Snack',
'Condiments',
'Additions'
);

CREATE TABLE IF NOT EXISTS items (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
merchant_id UUID NOT NULL REFERENCES merchants(id),
name VARCHAR NOT NULL,
product_category product_categories NOT NULL,
price INTEGER NOT NULL,
image_url VARCHAR NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at);
CREATE INDEX IF NOT EXISTS idx_items_product_category ON items(product_category);