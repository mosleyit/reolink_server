-- Drop triggers
DROP TRIGGER IF EXISTS update_camera_configs_updated_at ON camera_configs;
DROP TRIGGER IF EXISTS update_cameras_updated_at ON cameras;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS delete_old_events();

-- Drop tables (in reverse order of creation due to foreign keys)
DROP TABLE IF EXISTS camera_configs;
DROP TABLE IF EXISTS recordings;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS cameras;

-- Drop extensions (optional - comment out if other databases use them)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
-- DROP EXTENSION IF EXISTS timescaledb;

