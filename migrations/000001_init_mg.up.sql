BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS auth_user (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login TIMESTAMPTZ NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL,
    title VARCHAR(100) NOT NULL,
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id VARCHAR(50) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL,
    comment TEXT NOT NULL,
    post_id VARCHAR(50) REFERENCES posts(id) ON DELETE CASCADE,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

COMMIT;