-- +goose Up
CREATE TABLE goals (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID REFERENCES users(id) ON DELETE CASCADE,
    "year" INT NOT NULL,
    "month" INT NOT NULL,
    "income_goal" NUMERIC(12, 2) DEFAULT 0.00,
    "transaction_goal" INT DEFAULT 0,
    "estimated_average_sale_price" NUMERIC(12, 2) DEFAULT 0.00,
    "estimated_average_commission_rate" NUMERIC(5, 2) DEFAULT 0.00,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
