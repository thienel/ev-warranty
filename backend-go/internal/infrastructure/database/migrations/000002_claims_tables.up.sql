BEGIN;

CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'DRAFT',
    total_cost DECIMAL(15, 2) DEFAULT 0,
    approved_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_claims_deleted_at ON claims(deleted_at);
CREATE INDEX IF NOT EXISTS idx_claims_customer_id ON claims(customer_id);
CREATE INDEX IF NOT EXISTS idx_claims_vehicle_id ON claims(vehicle_id);

CREATE TABLE IF NOT EXISTS claim_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    part_category_id INTEGER NOT NULL,
    faulty_part_id UUID NOT NULL,
    replacement_part_id UUID,
    issue_description TEXT NOT NULL,
    status TEXT NOT NULL,
    type TEXT NOT NULL,
    cost DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_claim_items_claim FOREIGN KEY (claim_id)
    REFERENCES claims(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_claim_items_deleted_at ON claim_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_claim_items_claim_id ON claim_items(claim_id);

CREATE TABLE IF NOT EXISTS claim_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    attachment_type TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_claim_attachments_claim FOREIGN KEY (claim_id)
    REFERENCES claims(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_claim_attachments_deleted_at ON claim_attachments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_claim_attachments_claim_id ON claim_attachments(claim_id);

CREATE TABLE IF NOT EXISTS claim_histories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    status TEXT NOT NULL,
    changed_by UUID NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_claim_histories_claim FOREIGN KEY (claim_id)
    REFERENCES claims(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_claim_histories_deleted_at ON claim_histories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_claim_histories_claim_id ON claim_histories(claim_id);
CREATE INDEX IF NOT EXISTS idx_claim_histories_changed_by ON claim_histories(changed_by);

COMMIT;