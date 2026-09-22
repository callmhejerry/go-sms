CREATE TABLE IF NOT EXISTS scores (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id           UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject_id           UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    assessment_type_id   UUID NOT NULL REFERENCES assessment_types(id) ON DELETE CASCADE,
    class_arm_id         UUID NOT NULL REFERENCES class_arms(id) ON DELETE CASCADE,
    academic_session_id  UUID NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    score                NUMERIC(6, 2) NOT NULL CHECK (score >= 0),
    recorded_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT score_unique UNIQUE (tenant_id, student_id, subject_id, assessment_type_id, academic_session_id)
);

CREATE INDEX idx_scores_tenant_id ON scores(tenant_id);
CREATE INDEX idx_scores_student_id ON scores(student_id);
CREATE INDEX idx_scores_subject_id ON scores(subject_id);
CREATE INDEX idx_scores_class_arm ON scores(class_arm_id);
CREATE INDEX idx_scores_session ON scores(academic_session_id);