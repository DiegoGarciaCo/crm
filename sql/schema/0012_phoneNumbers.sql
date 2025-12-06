-- +goose Up
CREATE TABLE phone_numbers (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_id" UUID REFERENCES contacts(id) ON DELETE CASCADE,
    "phone_number" VARCHAR(20) NOT NULL,
    "type" VARCHAR(50) DEFAULT 'mobile',
    "is_primary" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
