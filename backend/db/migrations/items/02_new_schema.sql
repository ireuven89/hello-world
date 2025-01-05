-- +goose Up

create table if not exists bidders
(
    id varchar(36) primary key,
    name varchar(32) not null default ''
);

create table if not exists workers(
  id varchar(36) primary key,
  name varchar(32) not null default ''
);