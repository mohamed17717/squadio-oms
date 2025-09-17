CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Enums (ممكن تعملها CHECK لو مش عايز ENUM)
CREATE TYPE order_status AS ENUM (
    'draft',
    'pending_payment',
    'paid',
    'fulfillment_in_progress',
    'shipped',
    'completed',
    'cancelled'
);
CREATE TYPE payment_status AS ENUM (
    'pending',
    'authorized',
    'captured',
    'failed',
    'refunded',
    'partial_refunded'
);
CREATE TYPE refund_status AS ENUM ('pending', 'approved', 'rejected', 'processed');
-- Products & Variants
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    is_active BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku TEXT UNIQUE NOT NULL,
    attributes JSONB,
    -- { "color":"black", "size":"L" }
    price_minor INT NOT NULL CHECK (price_minor >= 0),
    currency CHAR(3) NOT NULL,
    is_active BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Inventory (بسيطة وواضحة)
CREATE TABLE inventory (
    variant_id UUID PRIMARY KEY REFERENCES product_variants(id) ON DELETE CASCADE,
    qty_on_hand INT NOT NULL DEFAULT 0 CHECK (qty_on_hand >= 0),
    qty_reserved INT NOT NULL DEFAULT 0 CHECK (qty_reserved >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Orders
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID REFERENCES customers(id) ON DELETE
    SET NULL,
        status order_status NOT NULL DEFAULT 'draft',
        subtotal_minor INT NOT NULL DEFAULT 0 CHECK (subtotal_minor >= 0),
        --   discount_minor INT NOT NULL DEFAULT 0 CHECK (discount_minor >= 0),
        --   shipping_minor INT NOT NULL DEFAULT 0 CHECK (shipping_minor >= 0),
        --   tax_minor INT NOT NULL DEFAULT 0 CHECK (tax_minor >= 0),
        total_minor INT NOT NULL DEFAULT 0 CHECK (total_minor >= 0),
        currency CHAR(3) NOT NULL,
        --   billing_address_id UUID REFERENCES addresses(id),
        --   shipping_address_id UUID REFERENCES addresses(id),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        version INT NOT NULL DEFAULT 1 -- optimistic locking
);
-- Order Items
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id UUID NOT NULL REFERENCES product_variants(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price_minor INT NOT NULL CHECK (unit_price_minor >= 0),
    currency CHAR(3) NOT NULL,
    tax_minor INT NOT NULL DEFAULT 0 CHECK (tax_minor >= 0),
    discount_minor INT NOT NULL DEFAULT 0 CHECK (discount_minor >= 0),
    line_total_minor INT NOT NULL CHECK (line_total_minor >= 0),
    UNIQUE(order_id, variant_id)
);
-- Payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    -- e.g. "stripe", "cod"
    status payment_status NOT NULL DEFAULT 'pending',
    amount_minor INT NOT NULL CHECK (amount_minor >= 0),
    currency CHAR(3) NOT NULL,
    external_ref TEXT,
    -- gateway payment_intent id
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Refunds
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    payment_id UUID REFERENCES payments(id) ON DELETE
    SET NULL,
        status refund_status NOT NULL DEFAULT 'pending',
        amount_minor INT NOT NULL CHECK (amount_minor >= 0),
        reason TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        processed_at TIMESTAMPTZ
);
-- Order Events (Audit لذيذ)
CREATE TABLE order_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    -- 'status_changed','item_added','payment_captured',...
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Indices مهمة
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_orders_open ON orders(status)
WHERE status IN (
        'pending_payment',
        'paid',
        'fulfillment_in_progress'
    );