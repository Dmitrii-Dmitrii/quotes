-- +goose Up
CREATE TABLE IF NOT EXISTS quotes (
    id UUID PRIMARY KEY,
    author TEXT NOT NULL,
    text TEXT NOT NULL
);