-- Enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Modify categories table
-- We assume categories table exists. If not, this might fail or we should create it.
-- But since this is an "alter" migration, we assume existence.
TRUNCATE TABLE categories CASCADE;

-- Drop constraints if they exist (to allow changing ID)
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_pkey;

-- Change ID to UUID
ALTER TABLE categories DROP COLUMN IF EXISTS id;
ALTER TABLE categories ADD COLUMN id UUID PRIMARY KEY DEFAULT uuid_generate_v4();

-- Add company_id
ALTER TABLE categories ADD COLUMN IF NOT EXISTS company_id UUID NOT NULL;

-- 2. Modify products table (fix company_id and category_id types)
TRUNCATE TABLE products CASCADE;
ALTER TABLE products DROP COLUMN IF EXISTS company_id;
ALTER TABLE products ADD COLUMN company_id UUID NOT NULL;
ALTER TABLE products DROP COLUMN IF EXISTS category_id;
ALTER TABLE products ADD COLUMN category_id UUID NOT NULL;

-- 3. Modify add_ons table
TRUNCATE TABLE add_ons CASCADE;
ALTER TABLE add_ons DROP COLUMN IF EXISTS company_id;
ALTER TABLE add_ons ADD COLUMN company_id UUID NOT NULL;

-- 4. Modify add_on_products table
TRUNCATE TABLE add_on_products CASCADE;
ALTER TABLE add_on_products DROP COLUMN IF EXISTS company_id;
ALTER TABLE add_on_products ADD COLUMN company_id UUID NOT NULL;

-- 5. Modify outlets table
TRUNCATE TABLE outlets CASCADE;
ALTER TABLE outlets DROP COLUMN IF EXISTS company_id;
ALTER TABLE outlets ADD COLUMN company_id UUID NOT NULL;

-- 6. Modify outlet_products table
TRUNCATE TABLE outlet_products CASCADE;
ALTER TABLE outlet_products DROP COLUMN IF EXISTS company_id;
ALTER TABLE outlet_products ADD COLUMN company_id UUID NOT NULL;

-- 7. Modify order_types table
TRUNCATE TABLE order_types CASCADE;
ALTER TABLE order_types DROP COLUMN IF EXISTS company_id;
ALTER TABLE order_types ADD COLUMN company_id UUID NOT NULL;

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
