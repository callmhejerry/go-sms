-- Classes (JSS1, JSS2)
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL, -- e.g JSS 1 etc
    level_order INT NOT NULL DEFAULT 0, -- for sorting
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT classes_tenant_id_name_key UNIQUE(tenant_id, name)
);

CREATE INDEX idx_classes_tenant_id ON classes(tenant_id);

-- ARMS ( A, B ARTS ETC)
CREATE TABLE IF NOT EXISTS class_arms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT class_arms_tenant_id_class_id_name_key UNIQUE(tenant_id, class_id, name)
);


CREATE INDEX idx_arms_tenant_id ON class_arms(tenant_id);
CREATE INDEX idx_arms_class_id ON class_arms(class_id);