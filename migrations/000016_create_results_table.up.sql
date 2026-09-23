CREATE TABLE IF NOT EXISTS results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    class_arm_id UUID NOT NULL REFERENCES class_arms(id) ON DELETE CASCADE,
    academic_session_id UUID NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    total_score          NUMERIC(8, 2) NOT NULL DEFAULT 0,
    max_total            NUMERIC(8, 2) NOT NULL DEFAULT 0,
    percentage           NUMERIC(5, 2) NOT NULL DEFAULT 0,
    grade                VARCHAR(5),
    remark               VARCHAR(50),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT result_unique UNIQUE (tenant_id, student_id, subject_id, academic_session_id)    
);

CREATE INDEX idx_results_tenant_id ON results(tenant_id);
CREATE INDEX idx_results_student_id ON results(student_id);
CREATE INDEX idx_results_class_arm ON results(class_arm_id);
CREATE INDEX idx_results_session ON results(academic_session_id);