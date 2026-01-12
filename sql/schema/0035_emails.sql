-- +goose Up
ALTER TABLE
    emails
ADD
    CONSTRAINT emails_email_unique UNIQUE (email_address);
