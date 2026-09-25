-- students
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    admission_number VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    gender VARCHAR(20) NOT NULL,
    date_of_birth DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    current_class_arm_id UUID REFERENCES class_arms(id) ON DELETE SET NULL,
    admission_session_id UUID REFERENCES academic_sessions(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT students_tenant_id_admission_number_key UNIQUE(tenant_id, admission_number)
);

CREATE INDEX idx_students_tenant_id ON students(id);
CREATE INDEX idx_students_admission_number ON students(tenant_id, admission_number);
CREATE INDEX idx_students_status ON students(tenant_id, status);
CREATE INDEX idx_students_current_class_arm ON students(current_class_arm_id);


-- Parents/Guardian
CREATE TABLE IF NOT EXISTS parents(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(225),
    phone_number VARCHAR(30) NOT NULL,
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_parents_tenant_id ON parents(tenant_id);
CREATE INDEX idx_parents_phone ON parents(tenant_id, phone_number);
CREATE INDEX idx_parents_email ON parents(tenant_id, email);

-- STUDENT <-> PARENT (MANY-MANY)
CREATE TABLE IF NOT EXISTS student_parents (
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    parent_id UUID NOT NULL REFERENCES parents(id) ON DELETE CASCADE,
    relationship VARCHAR(50) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (student_id, parent_id)
);

CREATE INDEX idx_student_parents_student_id ON student_parents(student_id);
CREATE INDEX idx_student_parents_parent_id ON student_parents(parent_id);