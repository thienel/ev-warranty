CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending');

CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   username VARCHAR(24) UNIQUE NOT NULL,
   email VARCHAR(100) UNIQUE NOT NULL,
   password_hash VARCHAR(255) NOT NULL,
   status user_status NOT NULL DEFAULT 'pending',
   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);