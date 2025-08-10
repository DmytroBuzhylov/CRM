CREATE TABLE IF NOT EXISTS recipe_items
(
    id              uuid PRIMARY KEY                  DEFAULT gen_random_uuid(),
    menu_item_id    uuid,
    ingredient_id   uuid,
    quantity_needed bigint,
    unit_of_measure varchar(255),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (menu_item_id) REFERENCES menu_items (id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients (id) ON DELETE CASCADE
);