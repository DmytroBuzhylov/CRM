CREATE TABLE IF NOT EXISTS notifications
(
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     uuid,
    user_id             uuid,

    type                varchar(50) NOT NULL CHECK (type IN ('UNSPECIFIED','ORDER', 'TASK', 'SYSTEM', 'CLIENT_REQUEST', 'SHIFT_STARTED', 'SHIFT_ENDED', 'SHIFT_OVERDUE', 'FOOD_READY', 'INCIDENT', 'EQUIPMENT_ALERT', 'SCREEN_MESSAGE', 'MANAGER_ALERT')),
    text                text,

    send_to_email       bool,
    send_to_app         bool,
    viewed              bool,

    resource_type       varchar(255) null CHECK ( type IN ('UNSPECIFIED', 'ORDER', 'TABLE', 'TASK', 'INVENTORY', 'MENU_ITEM', 'SHIFT', 'EMPLOYEE', 'KITCHEN_STATION', 'INCIDENT') ),
    related_resource_id uuid         null,

    created_at          bigint,

    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
