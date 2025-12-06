-- +goose Up
CREATE TABLE two_factor (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "secret" TEXT NOT NULL,
    "backupCodes" TEXT NOT NULL,
    "userId" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX two_factor_secret_idx ON two_factor(secret);

CREATE INDEX two_factor_userId_idx ON two_factor("userId");

-- +goose Down
DROP INDEX IF EXISTS two_factor_secret_idx;

DROP INDEX IF EXISTS two_factor_userId_idx;

DROP TABLE two_factor;
