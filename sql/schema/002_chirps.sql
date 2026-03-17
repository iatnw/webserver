-- +goose Up
create table chirps (
    id uuid primary key,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    body varchar(255) not null,
    user_id uuid not null references users(id) on delete cascade
);

-- +goose Down
drop table chirps;
