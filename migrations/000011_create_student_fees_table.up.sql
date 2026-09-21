CREATE TABLE IF NOT EXISTS student_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    fee_structure_id UUID NOT NULL REFERENCES fee_structures(id) ON DELETE RESTRICT,
    amount_kobo BIGINT NOT NULL CHECK (amount_kobo > 0),
    amount_paid_kobo BIGINT NOT NULL DEFAULT 0 CHECK (amount_paid_kobo >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
    due_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT student_fees_unique UNIQUE(tenant_id, student_id, fee_structure_id)
);

CREATE INDEX idx_student_fees_tenant_id ON student_fees(tenant_id);
CREATE INDEX idx_student_fees_student_id ON student_fees(student_id);
CREATE INDEX idx_student_fees_status ON student_fees(tenant_id, status);