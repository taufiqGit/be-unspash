ALTER TABLE products
DROP CONSTRAINT IF EXISTS products_unit_id_fkey;

DROP INDEX IF EXISTS idx_products_unit_id;

ALTER TABLE products
DROP COLUMN IF EXISTS unit_id;
