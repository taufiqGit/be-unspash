-- Drop junction tables first (dependent on discounts)
DROP TABLE IF EXISTS discount_order_types;
DROP TABLE IF EXISTS discount_target_products;
DROP TABLE IF EXISTS discount_target_categories;
DROP TABLE IF EXISTS discount_outlets;

-- Drop main table
DROP TABLE IF EXISTS discounts;

-- Drop enum types
DROP TYPE IF EXISTS discount_type;
DROP TYPE IF EXISTS discount_target;
DROP TYPE IF EXISTS discount_specific_target;
