-- +goose Up
CREATE TYPE appointment_type AS ENUM (
    'Listing-appointment',
    'Buyer-appointment',
    'no-type'
);

CREATE TYPE appointment_outcome AS ENUM (
    'no-outcome',
    'no',
    'yes',
    'rescheduled',
    'cancelled',
    'no-show'
);

CREATE TABLE appointments (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_id" UUID REFERENCES contacts(id) ON DELETE CASCADE,
    "assigned_to_id" UUID REFERENCES users(id) ON DELETE
    SET
        NULL,
        "title" VARCHAR(255) NOT NULL,
        "scheduled_at" TIMESTAMP WITH TIME ZONE NOT NULL,
        "location" VARCHAR(255) DEFAULT NULL,
        "type" appointment_type DEFAULT 'no-type',
        "outcome" appointment_outcome DEFAULT 'no-outcome',
        "note" TEXT DEFAULT NULL,
        "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
