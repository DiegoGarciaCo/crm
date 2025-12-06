-- +goose Up
CREATE TABLE invitation (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "organizationId" UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    "email" TEXT NOT NULL,
    "role" TEXT,
    "status" TEXT NOT NULL,
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "inviterId" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX invitation_organizationId_idx ON invitation("organizationId");

CREATE INDEX invitation_email_idx ON invitation(email);

-- +goose Down
DROP INDEX IF EXISTS invitation_organizationId_idx;

DROP INDEX IF EXISTS invitation_email_idx;

DROP TABLE invitation;
