ALTER TABLE users ADD COLUMN oauth_provider VARCHAR(32);
ALTER TABLE users ADD COLUMN oauth_id VARCHAR(64);
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

CREATE UNIQUE INDEX idx_users_oauth_provider_id
    ON users(oauth_provider, oauth_id)
    WHERE oauth_provider IS NOT NULL AND oauth_id IS NOT NULL;