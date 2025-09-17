CREATE TYPE office_type AS ENUM ('EVM', 'SC');

CREATE TABLE offices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    office_name VARCHAR(255) NOT NULL,
    office_type office_type NOT NULL,
    is_active BOOLEAN NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

