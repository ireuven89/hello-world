-- +goose Up
create table if not exists users
(
    id          varchar(255) primary key,
    uuid        char(36) not null,
    name        varchar(255),
    user_uuid   char(36),
    item        varchar(255),
    price       int,
    description varchar(255)
);