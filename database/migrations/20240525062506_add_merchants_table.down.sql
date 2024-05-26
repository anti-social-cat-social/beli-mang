DROP TABLE IF EXISTS merchants;

DROP INDEX IF EXISTS idx_merchants_created_at CASCADE;
DROP INDEX IF EXISTS idx_merchants_merchant_category CASCADE;

DROP TYPE IF EXISTS merchant_categories;