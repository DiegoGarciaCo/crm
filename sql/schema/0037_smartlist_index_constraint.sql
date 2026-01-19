-- +goose Up
ALTER TABLE
    smart_lists
ADD
    CONSTRAINT smart_lists_list_index_unique UNIQUE(list_index);
