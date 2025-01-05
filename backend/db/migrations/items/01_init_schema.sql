-- +goose Up

create table if not exists items
(
    id          varchar(255) primary key,
    uuid        varchar(36) not null,
    name        varchar(255),
    user_uuid   varchar(255),
    item        varchar(255),
    price       varchar(255),
    description varchar(255)
);
