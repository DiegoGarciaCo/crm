-- +goose Up
ALTER TABLE
    emails
ADD
    COLUMN is_verified BOOLEAN DEFAULT FALSE,
ADD
    COLUMN is_subscribed BOOLEAN DEFAULT FALSE;
