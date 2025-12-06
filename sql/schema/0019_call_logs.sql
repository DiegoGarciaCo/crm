-- +goose Up
CREATE TABLE contact_logs (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_id" UUID REFERENCES contacts(id) ON DELETE CASCADE,
    "contact_method" VARCHAR(50) NOT NULL,
    "created_by" UUID REFERENCES users(id) ON DELETE
    SET
        NULL,
        "note" TEXT DEFAULT NULL,
        "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
