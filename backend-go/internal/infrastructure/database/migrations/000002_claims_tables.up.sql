CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    total_cost DECIMAL(15, 2) DEFAULT 0,
    approved_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_claims_deleted_at ON claims(deleted_at);

CREATE TABLE IF NOT EXISTS claim_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    part_category_id INTEGER NOT NULL,
    faulty_part_id UUID NOT NULL,
    replacement_part_id UUID,
    issue_description TEXT NOT NULL,
    line_status TEXT NOT NULL,
    line_type TEXT NOT NULL,
    cost DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_claim_items_claim FOREIGN KEY (claim_id)
    REFERENCES claims(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_claim_items_claim_id ON claim_items(claim_id);
CREATE INDEX IF NOT EXISTS idx_claim_items_deleted_at ON claim_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_claim_items_part_category_id ON claim_items(part_category_id);
CREATE INDEX IF NOT EXISTS idx_claim_items_faulty_part_id ON claim_items(faulty_part_id);
CREATE INDEX IF NOT EXISTS idx_claim_items_line_status ON claim_items(line_status);

CREATE TABLE IF NOT EXISTS claim_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL,
    attachment_type TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_claim_attachments_claim FOREIGN KEY (claim_id)
    REFERENCES claims(id) ON DELETE CASCADE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_claim_attachments_deleted_at ON claim_attachments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_claim_attachments_claim_id ON claim_attachments(claim_id);
CREATE INDEX IF NOT EXISTS idx_claim_attachments_type ON claim_attachments(attachment_type);
