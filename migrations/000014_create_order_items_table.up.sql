CREATE TABLE IF NOT EXISTS order_items (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             order_id UUID NOT NULL REFERENCES orders(id),
                             menu_item_id UUID NOT NULL REFERENCES menu_items(id),
                             quantity INTEGER NOT NULL,
                             price_at_time_of_order DECIMAL(10, 2) NOT NULL,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                             updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);