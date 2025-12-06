-- +goose Up
CREATE TABLE contact_notes (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_id" UUID REFERENCES contacts(id) ON DELETE CASCADE,
    "note" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "created_by" UUID REFERENCES users(id) ON DELETE
    SET
        NULL
);
