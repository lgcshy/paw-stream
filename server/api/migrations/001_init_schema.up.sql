-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    nickname TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    disabled INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Create index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    owner_user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    location TEXT,
    publish_path TEXT UNIQUE NOT NULL,
    secret_hash TEXT NOT NULL,
    secret_cipher TEXT NOT NULL,
    secret_version INTEGER NOT NULL DEFAULT 1,
    disabled INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes on devices
CREATE INDEX IF NOT EXISTS idx_devices_owner ON devices(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_devices_path ON devices(publish_path);
