
-- +migrate Up
alter table user_sessions add column "deleted_at" timestamp default null;

-- +migrate Down
alter table user_sessions drop column "deleted_at";
