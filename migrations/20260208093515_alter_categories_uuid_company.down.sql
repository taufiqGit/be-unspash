-- Revert changes (Note: Data loss will occur as we cannot convert UUID back to INT reliably)

-- 1. Revert categories table
TRUNCATE TABLE categories CASCADE;
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_pkey;
ALTER TABLE categories DROP COLUMN id;
ALTER TABLE categories ADD COLUMN id SERIAL PRIMARY KEY;
ALTER TABLE categories DROP COLUMN IF EXISTS company_id;

-- 2. Revert products table
TRUNCATE TABLE products CASCADE;
ALTER TABLE products DROP COLUMN company_id;
ALTER TABLE products ADD COLUMN company_id INT NOT NULL;
ALTER TABLE products DROP COLUMN category_id;
ALTER TABLE products ADD COLUMN category_id INT NOT NULL;

-- 3. Revert add_ons table
TRUNCATE TABLE add_ons CASCADE;
ALTER TABLE add_ons DROP COLUMN company_id;
ALTER TABLE add_ons ADD COLUMN company_id INT NOT NULL;

-- 4. Revert add_on_products table
TRUNCATE TABLE add_on_products CASCADE;
ALTER TABLE add_on_products DROP COLUMN company_id;
ALTER TABLE add_on_products ADD COLUMN company_id INT NOT NULL;

-- 5. Revert outlets table
TRUNCATE TABLE outlets CASCADE;
ALTER TABLE outlets DROP COLUMN company_id;
ALTER TABLE outlets ADD COLUMN company_id INT NOT NULL;

-- 6. Revert outlet_products table
TRUNCATE TABLE outlet_products CASCADE;
ALTER TABLE outlet_products DROP COLUMN company_id;
ALTER TABLE outlet_products ADD COLUMN company_id INT NOT NULL;

-- 7. Revert order_types table
TRUNCATE TABLE order_types CASCADE;
ALTER TABLE order_types DROP COLUMN company_id;
ALTER TABLE order_types ADD COLUMN company_id INT NOT NULL;

-- Re-create indexes
DROP INDEX IF EXISTS idx_products_company;
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_add_ons_company;
DROP INDEX IF EXISTS idx_outlets_company;
DROP INDEX IF EXISTS idx_outlet_products_outlet;
DROP INDEX IF EXISTS idx_outlet_products_product;
DROP INDEX IF EXISTS idx_order_types_company;

CREATE INDEX idx_products_company ON products(company_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_add_ons_company ON add_ons(company_id);
CREATE INDEX idx_outlets_company ON outlets(company_id);
CREATE INDEX idx_order_types_company ON order_types(company_id);
