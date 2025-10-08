DROP INDEX IF EXISTS idx_claim_attachments_deleted_at;
DROP INDEX IF EXISTS idx_claim_attachments_type;
DROP INDEX IF EXISTS idx_claim_attachments_claim_id;
DROP INDEX IF EXISTS idx_claim_items_line_status;
DROP INDEX IF EXISTS idx_claim_items_faulty_part_id;
DROP INDEX IF EXISTS idx_claim_items_part_category_id;
DROP INDEX IF EXISTS idx_claim_items_deleted_at;
DROP INDEX IF EXISTS idx_claim_items_claim_id;
DROP INDEX IF EXISTS idx_claims_deleted_at;

DROP TABLE IF EXISTS claim_attachments CASCADE;
DROP TABLE IF EXISTS claim_items CASCADE;
DROP TABLE IF EXISTS claims CASCADE;