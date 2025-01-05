-- +goose Up

alter table users
 add column created_at timestamp not null default current_timestamp,
 add column updated_at timestamp not null default current_timestamp;
