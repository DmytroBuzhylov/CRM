CREATE TABLE IF NOT EXISTS tasks
(
    id              uuid PRIMARY KEY                  DEFAULT gen_random_uuid(),
    name            VARCHAR(255)             NOT NULL,
    description     TEXT,
    priority        INTEGER                  NOT NULL,
    status          VARCHAR(50)              NOT NULL,
    deadline        TIMESTAMP WITH TIME ZONE,
    assignee_id     uuid                     NULL,
    client_id       uuid                     NULL,
    organization_id uuid,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);