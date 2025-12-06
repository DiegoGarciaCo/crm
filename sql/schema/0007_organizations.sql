-- +goose Up
CREATE TABLE organization (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "slug" TEXT NOT NULL UNIQUE,
    "logo" TEXT,
    "createdAt" TIMESTAMPTZ NOT NULL,
    "metadata" TEXT
);

-- +goose Down
DROP TABLE organization;
