-- +goose Up
ALTER TABLE
    tags
ADD
    CONSTRAINT unique_user_tag UNIQUE (user_id, name);
