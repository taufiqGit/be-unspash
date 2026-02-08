-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id INT NOT NULL,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    unit VARCHAR(50),
    cost DECIMAL(15, 2) DEFAULT 0,
    price DECIMAL(15, 2) DEFAULT 0,
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create add_ons table
CREATE TABLE IF NOT EXISTS add_ons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(15, 2) DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create add_on_products table
CREATE TABLE IF NOT EXISTS add_on_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id INT NOT NULL,
    add_on_id UUID NOT NULL REFERENCES add_ons(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create outlets table
CREATE TABLE IF NOT EXISTS outlets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id INT NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    supervisor VARCHAR(255),
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create outlet_products table
CREATE TABLE IF NOT EXISTS outlet_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id INT NOT NULL,
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    stock INT DEFAULT 0,
    price DECIMAL(15, 2) DEFAULT 0,
    cost DECIMAL(15, 2) DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create order_types table
CREATE TABLE IF NOT EXISTS order_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_active_price_adjustment BOOLEAN DEFAULT FALSE,
    price_increase DECIMAL(15, 2) DEFAULT 0,
    price_decrease DECIMAL(15, 2) DEFAULT 0,
    increase_type VARCHAR(50), -- e.g., 'percentage' or 'fixed'
    decrease_type VARCHAR(50),
    increase_value DECIMAL(15, 2) DEFAULT 0,
    decrease_value DECIMAL(15, 2) DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_products_company ON products(company_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_add_ons_company ON add_ons(company_id);
CREATE INDEX idx_outlets_company ON outlets(company_id);
CREATE INDEX idx_outlet_products_outlet ON outlet_products(outlet_id);
CREATE INDEX idx_outlet_products_product ON outlet_products(product_id);
CREATE INDEX idx_order_types_company ON order_types(company_id);
