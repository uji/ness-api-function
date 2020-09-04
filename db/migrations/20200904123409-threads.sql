
-- +migrate Up
create extension pgcrypto;
create table if not exists threads(
  uuid uuid not null default gen_random_uuid(),
  title varchar not null,
  closed boolean not null default false,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp,
  constraint thread_pkey primary key (uuid)
);

-- +migrate Down
drop extension pgcrypto;
drop table if exists threads;
