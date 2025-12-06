-- +goose Up
CREATE TABLE deals (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_id" UUID REFERENCES contacts(id) ON DELETE CASCADE,
    "assigned_to_id" UUID REFERENCES users(id) ON DELETE
    SET
        NULL,
        "title" VARCHAR(255) NOT NULL,
        "price" INTEGER NOT NULL,
        "closing_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "earnest_money_due_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "mutual_acceptance_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "inspection_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "appraisal_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "final_walkthrough_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "possession_date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "commission" INTEGER DEFAULT NULL,
        "commission_split" INTEGER DEFAULT NULL,
        "property_address" VARCHAR(255) DEFAULT NULL,
        "property_city" VARCHAR(100) DEFAULT NULL,
        "property_state" VARCHAR(100) DEFAULT NULL,
        "property_zip_code" VARCHAR(20) DEFAULT NULL,
        "description" TEXT DEFAULT NULL,
        "stage_id" UUID REFERENCES stages(id) ON DELETE
    SET
        NULL,
        "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS deals;
