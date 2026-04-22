ALTER TABLE discounts
ADD COLUMN IF NOT EXISTS specific_target_type VARCHAR(20) NULL
CHECK (specific_target_type IN ('category', 'product'));

ALTER TABLE discounts
DROP COLUMN IF EXISTS priority;
