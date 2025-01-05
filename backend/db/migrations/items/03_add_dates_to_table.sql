-- +goose up

alter table bidders
    add column created_at timestamp default current_timestamp,
    add column updated_at timestamp default current_timestamp;

alter table items
    add column created_at timestamp default current_timestamp,
    add column updated_at timestamp default current_timestamp;

