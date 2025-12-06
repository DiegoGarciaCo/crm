-- +goose Up
CREATE TABLE "apikey" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" text,
    "start" text,
    "prefix" text,
    "key" text NOT NULL,
    "userId" UUID NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "refillInterval" integer,
    "refillAmount" integer,
    "lastRefillAt" timestamptz,
    "enabled" boolean,
    "rateLimitEnabled" boolean,
    "rateLimitTimeWindow" integer,
    "rateLimitMax" integer,
    "requestCount" integer,
    "remaining" integer,
    "lastRequest" timestamptz,
    "expiresAt" timestamptz,
    "createdAt" timestamptz NOT NULL,
    "updatedAt" timestamptz NOT NULL,
    "permissions" text,
    "metadata" text
);

CREATE INDEX "apikey_key_idx" ON "apikey" ("key");

CREATE INDEX "apikey_userId_idx" ON "apikey" ("userId");

-- +goose Down
DROP TABLE IF EXISTS "apikey";
