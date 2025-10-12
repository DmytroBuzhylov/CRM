CREATE TABLE IF NOT EXISTS orders (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        organization_id UUID NOT NULL REFERENCES organizations(id),
                        user_id UUID,
                        table_id UUID REFERENCES tables(id),
                        user_type VARCHAR(50) NOT NULL,
                        status VARCHAR(50) NOT NULL, -- 'pending', 'confirmed', 'completed', 'canceled'
                        notes TEXT,
                        total_price DECIMAL(10, 2) NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);