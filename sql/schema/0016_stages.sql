-- +goose Up
CREATE TYPE client_type AS ENUM ('buyer', 'seller');

CREATE TABLE stages (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT DEFAULT NULL,
    "client_type" client_type NOT NULL,
    "number_of_deals" INTEGER DEFAULT 0,
    "total_potential_income" INTEGER DEFAULT NULL,
    "order_index" INTEGER NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
