-- +goose Up
ALTER TABLE
    users
ADD
    COLUMN "trialAllowed" BOOLEAN DEFAULT TRUE NOT NULL;
