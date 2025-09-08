CREATE TYPE job_status AS ENUM ('Submitted', 'In Progress', 'Complete', 'Error', 'Cancelled');

CREATE TABLE IF NOT EXISTS extract_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status job_status DEFAULT 'Submitted',
    status_message TEXT DEFAULT 'Submitted, pending extraction',
    image_uuid UUID,
    message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_image FOREIGN KEY (image_uuid)
        REFERENCES images(id) ON DELETE CASCADE
);

CREATE INDEX idx_extract_jobs_image_uuid ON extract_jobs(image_uuid);

