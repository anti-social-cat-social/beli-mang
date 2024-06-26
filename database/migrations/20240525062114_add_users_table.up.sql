CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role
 AS ENUM (
'admin',
'user'
);

CREATE TABLE IF NOT EXISTS users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
role user_role NOT NULL,
username VARCHAR UNIQUE NOT NULL,
password VARCHAR NOT NULL,
email VARCHAR NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);