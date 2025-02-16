-- +goose Up

create table if not exists bidders
(
    id          varchar(255) primary key,
    uuid        varchar(255) not null,
    name        varchar(255) not null default '',
    description varchar(255) not null default '',
    item        varchar(255) not null default '',
    price       bigint       not null default 0,
    created_at  timestamp    not null default current_timestamp,
    updated_at  timestamp    not null default current_timestamp
);
