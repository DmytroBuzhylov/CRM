ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_organization
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id)
            ON DELETE SET NULL;