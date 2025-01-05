CREATE TABLE IF NOT EXISTS users (
    `id` varchar(36) primary key,
    `name` varchar(36),
);

CREATE TABLE IF NOT EXISTS players (
    `id` varchar(36) primary key,
    `name` varchar(36),
)

CREATE TABLE IF NOT EXISTS stands(
    `id` varchar(36) primary key,
    `name` varchar(10),
    `download_link` varchar(255),
)