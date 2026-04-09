CREATE TABLE IF NOT EXISTS cashier_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL,
    opening_cash NUMERIC(14,2) NOT NULL DEFAULT 0,
    closing_cash NUMERIC(14,2) NOT NULL DEFAULT 0,
    expected_cash NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cashier_shifts_company_id ON cashier_shifts(company_id);
CREATE INDEX IF NOT EXISTS idx_cashier_shifts_outlet_id ON cashier_shifts(outlet_id);
CREATE INDEX IF NOT EXISTS idx_cashier_shifts_user_id ON cashier_shifts(user_id);
CREATE INDEX IF NOT EXISTS idx_cashier_shifts_status ON cashier_shifts(status);
CREATE INDEX IF NOT EXISTS idx_cashier_shifts_start_time ON cashier_shifts(start_time);
