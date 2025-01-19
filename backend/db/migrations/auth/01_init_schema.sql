-- +goose Up

create table if not exists users
(
    id       varchar(255) primary key,
    user     varchar(255) not null default '' unique key,
    password varchar(255)
);
