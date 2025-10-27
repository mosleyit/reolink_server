-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_camera_groups_updated_at ON camera_groups;

-- Drop constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_role;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_cameras_group_id;
DROP INDEX IF EXISTS idx_camera_groups_name;

-- Remove group_id from cameras
ALTER TABLE cameras DROP COLUMN IF EXISTS group_id;

-- Drop tables
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS camera_groups;

