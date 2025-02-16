-- +goose up

alter table items
    add column link varchar(255),
    modify column price integer signed not null default 0;