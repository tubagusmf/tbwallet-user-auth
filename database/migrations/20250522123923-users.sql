
-- +migrate Up
create type roles as enum ('user', 'admin');

create table users (
    "id" serial primary key,
    "name" varchar(150) not null,
    "email" varchar(150) not null unique,
    "password_hash" text not null,
    "phone" varchar(20) not null,
    "kyc_status_id" int not null,
    "role" roles default 'user',
    "created_at" timestamp default current_timestamp,
    "updated_at" timestamp default current_timestamp,
    "deleted_at" timestamp default null
);

-- +migrate Down
drop table if exists users;
drop type if exists roles;
