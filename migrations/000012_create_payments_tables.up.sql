CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    amount_kobo BIGINT NOT NULL CHECK(amount_kobo > 0),
    payment_method TEXT NOT NULL DEFAULT 'cash',
    reference TEXT,
    received_by UUID REFERENCES users(id) ON DELETE SET NULL,
    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX idx_payments_student_id ON payments(student_id);


CREATE TABLE IF NOT EXISTS payment_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    student_fee_id UUID NOT NULL REFERENCES student_fees(id) ON DELETE RESTRICT,
    amount_kobo BIGINT NOT NULL CHECK (amount_kobo > 0),

    CONSTRAINT payment_allocation_unique UNIQUE (payment_id, student_fee_id) 
);

CREATE INDEX idx_payment_allocations_payment_id ON payment_allocations(payment_id);
CREATE INDEX idx_payment_allocations_student_fee_id ON payment_allocations(student_fee_id);