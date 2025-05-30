
-- +migrate Up
create type status as enum ('pending', 'approved', 'rejected');

create table kyc_documents (
    "id" serial primary key,
    "user_id" int not null references users("id") on delete cascade,
    "document_type" varchar(50) not null,
    "document_url" text not null,
    "status" status default 'pending',
    "created_at" timestamp default current_timestamp,
    "updated_at" timestamp default current_timestamp
);

-- +migrate Down
drop table if exists kyc_documents;
drop type if exists status;
