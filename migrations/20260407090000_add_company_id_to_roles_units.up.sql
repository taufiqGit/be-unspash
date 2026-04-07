ALTER TABLE roles ADD COLUMN IF NOT EXISTS company_id UUID;
UPDATE roles r
SET company_id = u.company_id
FROM users u
WHERE r.company_id IS NULL
  AND (r.created_by = u.id OR r.updated_by = u.id);
ALTER TABLE roles ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_company_id_fkey;
ALTER TABLE roles ADD CONSTRAINT roles_company_id_fkey FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS idx_roles_company_id ON roles(company_id);

ALTER TABLE units ADD COLUMN IF NOT EXISTS company_id UUID;
UPDATE units uo
SET company_id = u.company_id
FROM users u
WHERE uo.company_id IS NULL
  AND (uo.created_by = u.id OR uo.updated_by = u.id);
ALTER TABLE units ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE units DROP CONSTRAINT IF EXISTS units_company_id_fkey;
ALTER TABLE units ADD CONSTRAINT units_company_id_fkey FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS idx_units_company_id ON units(company_id);
