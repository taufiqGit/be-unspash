ALTER TABLE add_on_products
DROP COLUMN IF EXISTS product_id;

DROP INDEX IF EXISTS idx_add_on_products_product_id;
