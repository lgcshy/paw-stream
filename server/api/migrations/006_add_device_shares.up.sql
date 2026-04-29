CREATE TABLE IF NOT EXISTS device_shares (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL,
    shared_by_user_id TEXT NOT NULL,
    shared_to_user_id TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_by_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_to_user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(device_id, shared_to_user_id)
);

CREATE INDEX idx_device_shares_device ON device_shares(device_id);
CREATE INDEX idx_device_shares_to ON device_shares(shared_to_user_id);
