
-- +migrate Up
UPDATE users
SET kyc_status_id = NULL
WHERE kyc_status_id IS NOT NULL
  AND kyc_status_id NOT IN (SELECT id FROM kyc_documents);

-- +migrate Down
UPDATE users
SET kyc_status_id = 1
WHERE kyc_status_id IS NULL;
