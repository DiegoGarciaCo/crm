-- +goose Up
ALTER TABLE
    smart_lists
ALTER COLUMN
    filter_criteria
SET
    DEFAULT '{}'::jsonb;
