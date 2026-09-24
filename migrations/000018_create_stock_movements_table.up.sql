CREATE TABLE IF NOT EXISTS stock_movements(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('in', 'out', 'adjust')),
    quantity INT NOT NULL CHECK (quantity > 0),
    reason TEXT,
    reference TEXT,
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



CREATE INDEX idx_stock_movements_tenant_id ON stock_movements(tenant_id);
CREATE INDEX idx_stock_movements_item_id ON stock_movements(item_id);
CREATE INDEX idx_stock_movements_type ON stock_movements(movement_type);