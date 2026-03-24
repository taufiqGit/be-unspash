ALTER TABLE add_on_products
ADD COLUMN IF NOT EXISTS product_id UUID REFERENCES products(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_add_on_products_product_id ON add_on_products(product_id);
