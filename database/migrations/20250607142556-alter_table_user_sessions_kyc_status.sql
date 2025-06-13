-- +migrate Up
alter table user_sessions add column "kyc_status_id" int not null references kyc_documents("id") on delete cascade;

-- +migrate Down
alter table user_sessions drop column "kyc_status_id";