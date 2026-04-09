DROP INDEX IF EXISTS idx_stock_movements_created_at;
DROP INDEX IF EXISTS idx_stock_movements_reference_id;
DROP INDEX IF EXISTS idx_stock_movements_reference_type;
DROP INDEX IF EXISTS idx_stock_movements_type;
DROP INDEX IF EXISTS idx_stock_movements_outlet_id;
DROP INDEX IF EXISTS idx_stock_movements_product_id;
DROP TABLE IF EXISTS stock_movements;

DROP INDEX IF EXISTS idx_stocks_outlet_id;
DROP INDEX IF EXISTS idx_stocks_product_id;
DROP TABLE IF EXISTS stocks;

DROP INDEX IF EXISTS idx_purchase_details_product_id;
DROP INDEX IF EXISTS idx_purchase_details_purchase_id;
DROP TABLE IF EXISTS purchase_details;

DROP INDEX IF EXISTS idx_purchases_created_at;
DROP INDEX IF EXISTS idx_purchases_status;
DROP INDEX IF EXISTS idx_purchases_outlet_id;
DROP INDEX IF EXISTS idx_purchases_user_id;
DROP INDEX IF EXISTS idx_purchases_company_id;
DROP TABLE IF EXISTS purchases;
