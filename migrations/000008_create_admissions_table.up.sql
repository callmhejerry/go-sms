CREATE TABLE IF NOT EXISTS admissions (
    id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    academic_session_id UUID NOT NULL REFERENCES academic_sessions(id),

    -- Applicant information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    gender VARCHAR(20) NOT NULL,
    date_of_birth DATE NOT NULL,

    -- preferred class
    preferred_class_id UUID REFERENCES classes(id) ON DELETE SET NULL,

    -- status workflow
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    -- pending --> under_review --> accepted/rejected

    admission_number VARCHAR(50),
    student_id UUID REFERENCES students(id) ON DELETE SET NULL,

    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    rejection_reason  TEXT,

    -- Simple parent info (we can expand later)
    parent_first_name VARCHAR(200) NOT NULL,
    parent_last_name VARCHAR(200) NOT NULL,
    parent_phone_number  VARCHAR(30) NOT NULL,
    parent_email     VARCHAR(255) NOT NULL,
    parent_relationship VARCHAR(100) NOT NULL,
    
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admissions_tenant_id ON admissions(tenant_id);
CREATE INDEX idx_admissions_status ON admissions(tenant_id, status);
CREATE INDEX idx_admissions_session ON admissions(academic_session_id);