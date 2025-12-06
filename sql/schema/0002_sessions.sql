-- +goose Up
CREATE TABLE SESSION (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "token" TEXT NOT NULL UNIQUE,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ NOT NULL,
    "ipAddress" TEXT,
    "userAgent" TEXT,
    "userId" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "impersonatedBy" UUID,
    "activeOrganizationId" UUID
);

CREATE INDEX sessions_userId_idx ON SESSION("userId");

-- +goose Down
DROP INDEX IF EXISTS sessions_userId_idx;

DROP TABLE SESSION;
