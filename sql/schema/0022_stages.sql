-- +goose Up
ALTER TABLE
    stages
ADD
    COLUMN owner_id UUID REFERENCES users(id) ON DELETE CASCADE;
