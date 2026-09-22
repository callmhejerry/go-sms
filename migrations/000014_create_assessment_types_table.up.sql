CREATE TABLE IF NOT EXISTS assessment_types(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    max_score NUMERIC(6, 2) NOT NULL CHECK (max_score > 0),
    weight NUMERIC(5, 2) DEFAULT 0,
    is_exam BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT assessment_type_unique UNIQUE(tenant_id, name)
);
CREATE INDEX idx_assessment_types_tenant_id ON assessment_types(tenant_id);