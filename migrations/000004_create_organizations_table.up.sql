CREATE TABLE IF NOT EXISTS organizations
(
    id            uuid PRIMARY KEY         DEFAULT gen_random_uuid(),
    name          varchar(255),
    description   text,
    owner_user_id uuid,
    created_at    timestamp with time zone default current_timestamp,
    FOREIGN KEY (owner_user_id) REFERENCES users (id)
);


CREATE TABLE IF NOT EXISTS organizations_users
(
    organization_id uuid,
    user_id         uuid,
    role            varchar(50),
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);