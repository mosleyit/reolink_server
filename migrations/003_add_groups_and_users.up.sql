-- Create camera_groups table
CREATE TABLE IF NOT EXISTS camera_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_camera_groups_name ON camera_groups(name);

-- Add group_id to cameras table
ALTER TABLE cameras
    ADD COLUMN IF NOT EXISTS group_id UUID REFERENCES camera_groups(id) ON DELETE SET NULL;

CREATE INDEX idx_cameras_group_id ON cameras(group_id);

-- Create users table for API authentication
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Add check constraint for user roles
ALTER TABLE users
    ADD CONSTRAINT check_user_role
    CHECK (role IN ('admin', 'user', 'viewer'));

-- Apply updated_at trigger to camera_groups
CREATE TRIGGER update_camera_groups_updated_at
    BEFORE UPDATE ON camera_groups
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply updated_at trigger to users
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create default admin user (password: 'admin' - CHANGE THIS!)
-- Password hash is bcrypt hash of 'admin'
INSERT INTO users (username, password_hash, email, role)
VALUES (
    'admin',
    '$2a$10$63YFx5c6CA8YorjpQYgvyuuPhDg39X5kYygxuo2HlHd8Z5EtYBmp2',
    'admin@localhost',
    'admin'
) ON CONFLICT (username) DO NOTHING;

