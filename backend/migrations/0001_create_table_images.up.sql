CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path VARCHAR(255),
    mimetype VARCHAR(255),
    size INT,
    encoded BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP Default NOW()
);

