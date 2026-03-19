-- ============================================================
-- 1. Tabel utama: discounts
-- ============================================================
CREATE TABLE IF NOT EXISTS discounts (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id           UUID NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    name                 VARCHAR(255) NOT NULL,

    -- 4 jenis diskon: product_rp | product_pct | receipt_rp | receipt_pct
    type                 VARCHAR(50)  NOT NULL
                             CHECK (type IN ('product_rp', 'product_pct', 'receipt_rp', 'receipt_pct')),

    -- Nilai diskon: nominal Rp atau persentase (%)
    discount_value       NUMERIC(15, 2) NOT NULL DEFAULT 0,

    -- Hanya untuk tipe product_pct & receipt_pct
    max_amount           NUMERIC(15, 2) NULL,

    -- Hanya untuk tipe receipt_rp & receipt_pct
    min_purchase         NUMERIC(15, 2) NULL,

    -- Target diskon: 'all' | 'specific' — hanya untuk tipe product_rp & product_pct
    target_type          VARCHAR(20)  NULL
                             CHECK (target_type IN ('all', 'specific')),

    -- Sub-target saat target_type = 'specific': 'category' | 'product'
    specific_target_type VARCHAR(20)  NULL
                             CHECK (specific_target_type IN ('category', 'product')),

    -- Toggle "Terapkan ke Jenis Pesanan Tertentu"
    apply_to_order_types BOOLEAN NOT NULL DEFAULT FALSE,

    created_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_discounts_company_id ON discounts(company_id);
CREATE INDEX idx_discounts_type       ON discounts(type);


-- ============================================================
-- 2. Junction table: discount ↔ outlets  (many-to-many)
--    Satu diskon bisa berlaku di banyak cabang
-- ============================================================
CREATE TABLE IF NOT EXISTS discount_outlets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    outlet_id   UUID NOT NULL REFERENCES outlets(id)   ON DELETE CASCADE,

    CONSTRAINT uq_discount_outlet UNIQUE (discount_id, outlet_id)
);

CREATE INDEX idx_discount_outlets_discount_id ON discount_outlets(discount_id);
CREATE INDEX idx_discount_outlets_outlet_id   ON discount_outlets(outlet_id);


-- ============================================================
-- 3. Junction table: discount ↔ categories  (many-to-many)
--    Aktif saat specific_target_type = 'category'
-- ============================================================
CREATE TABLE IF NOT EXISTS discount_target_categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id)   ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id)  ON DELETE CASCADE,

    CONSTRAINT uq_discount_category UNIQUE (discount_id, category_id)
);

CREATE INDEX idx_discount_target_categories_discount_id  ON discount_target_categories(discount_id);
CREATE INDEX idx_discount_target_categories_category_id  ON discount_target_categories(category_id);


-- ============================================================
-- 4. Junction table: discount ↔ products  (many-to-many)
--    Aktif saat specific_target_type = 'product'
-- ============================================================
CREATE TABLE IF NOT EXISTS discount_target_products (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id UUID NOT NULL REFERENCES discounts(id)  ON DELETE CASCADE,
    product_id  UUID NOT NULL REFERENCES products(id)   ON DELETE CASCADE,

    CONSTRAINT uq_discount_product UNIQUE (discount_id, product_id)
);

CREATE INDEX idx_discount_target_products_discount_id ON discount_target_products(discount_id);
CREATE INDEX idx_discount_target_products_product_id  ON discount_target_products(product_id);


-- ============================================================
-- 5. Junction table: discount ↔ order_types  (many-to-many)
--    Aktif saat apply_to_order_types = TRUE
-- ============================================================
CREATE TABLE IF NOT EXISTS discount_order_types (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discount_id   UUID NOT NULL REFERENCES discounts(id)    ON DELETE CASCADE,
    order_type_id UUID NOT NULL REFERENCES order_types(id)  ON DELETE CASCADE,

    CONSTRAINT uq_discount_order_type UNIQUE (discount_id, order_type_id)
);

CREATE INDEX idx_discount_order_types_discount_id    ON discount_order_types(discount_id);
CREATE INDEX idx_discount_order_types_order_type_id  ON discount_order_types(order_type_id);
