
CREATE TABLE IF NOT EXISTS embed_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status job_status DEFAULT 'Submitted',
    status_message TEXT,
    image_uuid UUID,
    message TEXT,
    embedded_uuid UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_image FOREIGN KEY (image_uuid)
        REFERENCES images(id) ON DELETE CASCADE,
    CONSTRAINT fk_embed_image FOREIGN KEY (embedded_uuid)
        REFERENCES images(id) ON DELETE CASCADE
);

CREATE INDEX idx_embed_jobs_image_uuid ON embed_jobs(image_uuid);
CREATE INDEX idx_embed_jobs_embedded_uuid ON embed_jobs(embedded_uuid);

