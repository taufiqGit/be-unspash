ALTER TABLE products
ADD COLUMN IF NOT EXISTS unit_id UUID;

CREATE INDEX IF NOT EXISTS idx_products_unit_id ON products(unit_id);

ALTER TABLE products
DROP CONSTRAINT IF EXISTS products_unit_id_fkey;

ALTER TABLE products
ADD CONSTRAINT products_unit_id_fkey
FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL;
