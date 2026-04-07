DROP INDEX IF EXISTS idx_recipes_is_active;
DROP INDEX IF EXISTS idx_recipes_ingredient_id;
DROP INDEX IF EXISTS idx_recipes_product_id;
DROP INDEX IF EXISTS idx_recipes_company_id;
DROP TABLE IF EXISTS recipes;

DROP INDEX IF EXISTS idx_suppliers_is_active;
DROP INDEX IF EXISTS idx_suppliers_name;
DROP INDEX IF EXISTS idx_suppliers_company_id;
DROP TABLE IF EXISTS suppliers;
