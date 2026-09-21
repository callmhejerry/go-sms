-- Fee types (Tuition, Bus fee, etc)
CREATE TABLE IF NOT EXISTS fee_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_optional BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tenant_id_name_unique UNIQUE (tenant_id, name)   
);

CREATE INDEX idx_fee_types_tenant_id ON fee_types(tenant_id);

-- Fee Structures (amount per session + class)
CREATE TABLE IF NOT EXISTS fee_structures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    fee_type_id UUID NOT NULL REFERENCES fee_types(id) ON DELETE CASCADE,
    academic_session_id UUID NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE, -- NULL APPLIES TO ALL CLASSES
    amount_kobo BIGINT NOT NULL CHECK(amount_kobo > 0),
    due_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fee_structures_unique UNIQUE(tenant_id, fee_type_id, academic_session_id, class_id)
);


CREATE INDEX idx_fee_structures_tenant_id ON fee_structures(tenant_id);
CREATE INDEX idx_fee_structures_session ON fee_structures(academic_session_id);
CREATE INDEX idx_fee_structures_class ON fee_structures(class_id);