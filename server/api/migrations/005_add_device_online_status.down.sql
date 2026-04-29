-- SQLite doesn't support DROP COLUMN before 3.35.0, recreate table
CREATE TABLE devices_backup AS SELECT id, owner_user_id, name, location, publish_path, secret_hash, secret_cipher, secret_version, disabled, created_at, updated_at FROM devices;
DROP TABLE devices;
ALTER TABLE devices_backup RENAME TO devices;
CREATE INDEX idx_devices_owner ON devices(owner_user_id);
CREATE INDEX idx_devices_path ON devices(publish_path);
