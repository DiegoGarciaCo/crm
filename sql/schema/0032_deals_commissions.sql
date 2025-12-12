-- +goose Up
ALTER TABLE
    deals
ALTER COLUMN
    commission TYPE numeric(10, 2) USING commission::numeric,
ALTER COLUMN
    commission_split TYPE numeric(10, 2) USING commission_split::numeric;
