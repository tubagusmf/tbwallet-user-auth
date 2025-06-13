-- +migrate Up
alter table users ADD CONSTRAINT fk_kyc_status_id FOREIGN KEY ("kyc_status_id") REFERENCES kyc_documents("id") ON DELETE CASCADE;

-- +migrate Down
alter table users DROP CONSTRAINT fk_kyc_status_id;