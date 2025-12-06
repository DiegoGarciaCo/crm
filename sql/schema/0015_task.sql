-- +goose Up
-- Create custom ENUM types first
CREATE TYPE task_type AS ENUM (
    'call',
    'email',
    'follow-up',
    'text',
    'showing',
    'closing',
    'open-house',
    'thank-you'
);

CREATE TYPE task_status AS ENUM (
    'pending',
    'completed',
    'cancelled'
);

CREATE TYPE task_priority AS ENUM ('low', 'normal', 'high');

-- Now create the tasks table
CREATE TABLE tasks (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_id" UUID REFERENCES contacts(id) ON DELETE CASCADE,
    "assigned_to_id" UUID REFERENCES users(id) ON DELETE
    SET
        NULL,
        "title" VARCHAR(255) NOT NULL,
        "type" task_type DEFAULT 'follow-up',
        "date" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
        "status" task_status DEFAULT 'pending',
        "priority" task_priority DEFAULT 'normal',
        "note" TEXT DEFAULT NULL,
        "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS tasks;

DROP TYPE IF EXISTS task_type;

DROP TYPE IF EXISTS task_status;

DROP TYPE IF EXISTS task_priority;
