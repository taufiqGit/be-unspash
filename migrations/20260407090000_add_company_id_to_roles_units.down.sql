DROP INDEX IF EXISTS idx_units_company_id;
ALTER TABLE units DROP CONSTRAINT IF EXISTS units_company_id_fkey;
ALTER TABLE units DROP COLUMN IF EXISTS company_id;

DROP INDEX IF EXISTS idx_roles_company_id;
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_company_id_fkey;
ALTER TABLE roles DROP COLUMN IF EXISTS company_id;
