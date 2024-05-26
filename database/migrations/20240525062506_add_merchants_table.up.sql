CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE merchant_categories
 AS ENUM (
'SmallRestaurant',
'MediumRestaurant',
'LargeRestaurant',
'MerchandiseRestaurant',
'BoothKiosk',
'ConvenienceStore'
);

CREATE TABLE IF NOT EXISTS merchants (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name VARCHAR NOT NULL,
merchant_category merchant_categories NOT NULL,
image_url VARCHAR NOT NULL,
location_lat REAL NOT NULL,
location_long REAL NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_merchants_created_at ON merchants(created_at);
CREATE INDEX IF NOT EXISTS idx_merchants_merchant_category ON merchants(merchant_category);