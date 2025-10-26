-- app/backend/service-b/db/migrations/001_create_templates_table.sql
CREATE SCHEMA IF NOT EXISTS templates;



-- Users table
CREATE TABLE IF NOT EXISTS templates.userstemps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX idx_users_email ON templates.userstemps(email);
CREATE INDEX idx_users_username ON templates.userstemps(username);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON templates.userstemps
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();