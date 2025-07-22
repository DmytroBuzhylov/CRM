CREATE TABLE IF NOT EXISTS tasks (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       description TEXT,
                       priority INTEGER NOT NULL,
                       status VARCHAR(50) NOT NULL,
                       deadline TIMESTAMP WITH TIME ZONE,
                       assignee_id bigint,
                       client_id bigint,
                       organization_id bigint,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);