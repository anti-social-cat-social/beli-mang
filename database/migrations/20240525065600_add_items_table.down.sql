DROP TABLE IF EXISTS items;

DROP INDEX IF EXISTS idx_items_created_at CASCADE;
DROP INDEX IF EXISTS idx_items_product_category CASCADE;

DROP TYPE IF EXISTS product_categories;