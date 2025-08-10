CREATE TABLE IF NOT EXISTS ingredients
(
    id               uuid PRIMARY KEY         DEFAULT gen_random_uuid(),
    organization_id  uuid,
    name             VARCHAR(255) NOT NULL,
    unit             VARCHAR(255) NOT NULL,
    quantity         bigint,
    minimum_quantity bigint,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,


    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ingredients_name ON ingredients (name);
CREATE INDEX IF NOT EXISTS idx_ingredients_id ON ingredients (id);