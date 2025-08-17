CREATE TABLE IF NOT EXISTS unit_conversions
(
    id         uuid PRIMARY KEY         DEFAULT gen_random_uuid(),
    from_unit  VARCHAR(255) NOT NULL,
    to_unit    VARCHAR(255) NOT NULL,
    multiplier NUMERIC(10, 4) NOT NULL,

    CONSTRAINT unique_conversion UNIQUE (from_unit, to_unit)
);

INSERT INTO unit_conversions (from_unit, to_unit, multiplier)
VALUES
    ('kg', 'g', 1000),
    ('g', 'kg', 0.001),
    ('l', 'ml', 1000),
    ('ml', 'l', 0.001),
    ('kg', 'kg', 1),
    ('g', 'g', 1),
    ('l', 'l', 1),
    ('ml', 'ml', 1);

CREATE INDEX IF NOT EXISTS idx_conversions_from_to ON unit_conversions (from_unit, to_unit);