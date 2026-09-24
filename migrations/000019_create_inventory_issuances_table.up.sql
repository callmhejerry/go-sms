CREATE TABLE IF NOT EXISTS inventory_issuances (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    item_id              UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    quantity             INT NOT NULL CHECK (quantity > 0),
    issued_to_type       VARCHAR(30) NOT NULL,          -- student, department, staff
    issued_to_id         UUID NOT NULL,
    issued_by            UUID REFERENCES users(id) ON DELETE SET NULL,
    academic_session_id  UUID REFERENCES academic_sessions(id) ON DELETE SET NULL,
    notes                TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inventory_issuances_tenant_id ON inventory_issuances(tenant_id);
CREATE INDEX idx_inventory_issuances_item_id ON inventory_issuances(item_id);
CREATE INDEX idx_inventory_issuances_issued_to ON inventory_issuances(issued_to_type, issued_to_id);