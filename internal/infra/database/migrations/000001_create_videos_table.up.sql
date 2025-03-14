CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_path VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    upload_status VARCHAR(50) NOT NULL DEFAULT 'none',
    error_message TEXT,
    hls_path VARCHAR(255),
    manifest_path VARCHAR(255),
    s3_url VARCHAR(255),
    s3_manifest_url VARCHAR(255),
    segment_key VARCHAR(255),
    manifest_key VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
); 