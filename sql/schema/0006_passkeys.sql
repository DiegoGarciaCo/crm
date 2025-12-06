-- +goose Up
CREATE TABLE passkey (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" TEXT,
    "publicKey" TEXT NOT NULL,
    "userId" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "credentialID" TEXT NOT NULL,
    "counter" INTEGER NOT NULL,
    "deviceType" TEXT NOT NULL,
    "backedUp" BOOLEAN NOT NULL,
    "transports" TEXT,
    "createdAt" TIMESTAMPTZ,
    "aaguid" TEXT
);

CREATE INDEX passkey_userId_idx ON passkey("userId");

CREATE INDEX passkey_credentialID_idx ON passkey("credentialID");

-- +goose Down
DROP INDEX IF EXISTS passkey_userId_idx;

DROP INDEX IF EXISTS passkey_credentialID_idx;

DROP TABLE passkey;
