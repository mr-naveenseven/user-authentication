CREATE TABLE user_accounts (
    id SERIAL PRIMARY KEY,                  -- auto-incrementing primary key
    username VARCHAR(50) UNIQUE NOT NULL,   -- unique username
    email VARCHAR(100) UNIQUE NOT NULL,     -- unique email
    password_hash VARCHAR(255) NOT NULL,    -- hashed password
    is_active BOOLEAN DEFAULT TRUE,         -- account active status
    is_locked BOOLEAN DEFAULT FALSE,        -- account locked status
    password_modified_at TIMESTAMP,         -- timestamp of last password modification
    -- metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- record creation timestamp
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- record modification timestamp
);

