
-- +migrate Up
create table user_sessions (
    "id" serial primary key,
    "user_id" int not null references users("id") on delete cascade,
    "token" text not null unique,
    "expires_at" timestamp not null,
    "created_at" timestamp default current_timestamp
);

-- +migrate Down
drop table if exists user_sessions;
