-- +goose Up
CREATE TABLE member (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "organizationId" UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    "userId" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "role" TEXT NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL
);

CREATE INDEX member_organizationId_idx ON member("organizationId");

CREATE INDEX member_userId_idx ON member("userId");

-- +goose Down
DROP INDEX IF EXISTS member_organizationId_idx;

DROP INDEX IF EXISTS member_userId_idx;

DROP TABLE member;
