-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    phone      VARCHAR(50),
    email      VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_customers_company_id ON customers(company_id);
CREATE INDEX idx_customers_name ON customers(name);
