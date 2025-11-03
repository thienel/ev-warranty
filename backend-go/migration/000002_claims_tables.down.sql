DROP INDEX IF EXISTS idx_claim_histories_changed_by;
DROP INDEX IF EXISTS idx_claim_histories_claim_id;
DROP INDEX IF EXISTS idx_claim_attachments_claim_id;
DROP INDEX IF EXISTS idx_claim_histories_deleted_at;
DROP INDEX IF EXISTS idx_claim_attachments_deleted_at;
DROP INDEX IF EXISTS idx_claim_items_claim_id;
DROP INDEX IF EXISTS idx_claim_items_deleted_at;
DROP INDEX IF EXISTS idx_claims_vehicle_id;
DROP INDEX IF EXISTS idx_claims_customer_id;
DROP INDEX IF EXISTS idx_claims_deleted_at;

DROP TABLE IF EXISTS claim_histories CASCADE;
DROP TABLE IF EXISTS claim_attachments CASCADE;
DROP TABLE IF EXISTS claim_items CASCADE;
DROP TABLE IF EXISTS claims CASCADE;