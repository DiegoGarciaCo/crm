-- +goose Up
ALTER TABLE
    smart_lists
ADD
    COLUMN list_index SERIAL NOT NULL;
