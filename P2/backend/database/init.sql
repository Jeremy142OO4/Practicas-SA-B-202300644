BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS roles (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT roles_name_not_blank CHECK (BTRIM(name) <> '')
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name_encrypted TEXT NOT NULL,
    email_encrypted TEXT NOT NULL,
    email_hash CHAR(64) NOT NULL UNIQUE,
    password_encrypted TEXT NOT NULL,
    role_id SMALLINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_fk
        FOREIGN KEY (role_id) REFERENCES roles(id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT users_status_allowed
        CHECK (status IN ('ACTIVE', 'INACTIVE', 'BLOCKED')),
    CONSTRAINT users_email_hash_format
        CHECK (email_hash ~ '^[0-9a-f]{64}$'),
    CONSTRAINT users_encrypted_fields_not_blank
        CHECK (
            BTRIM(name_encrypted) <> '' AND
            BTRIM(email_encrypted) <> '' AND
            BTRIM(password_encrypted) <> ''
        )
);

CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS roles_set_updated_at ON roles;
CREATE TRIGGER roles_set_updated_at
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS users_set_updated_at ON users;
CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

INSERT INTO roles (name, description)
VALUES
    ('ADMIN', 'Administrator with access to all protected routes'),
    ('CLIENTE', 'Client with access to shared protected routes')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;

-- AES-256-GCM seed values use the development AES_KEY from backend/.env.example.
-- email_hash uses HMAC-SHA256 with EMAIL_INDEX_KEY from backend/.env.example.
INSERT INTO users (
    id,
    name_encrypted,
    email_encrypted,
    email_hash,
    password_encrypted,
    role_id,
    status
)
VALUES
    (
        '00000000-0000-4000-8000-000000000001',
        '9Bq14xp/IkIuSpMRBJ47/86Aznwv9diso1DOHjXK5PGdxxpzn0Lvh1xIVF/9GC+ecQ==',
        'I0Y9WZVwurC1xr0F3IpkY7ykbr7PUEnJaXdNtdwedOi++QZEefR6L+ZmCJh/RHhL',
        'abf0788e2783e4a0caf57c4725c4488dbba89ff1891cc5823d8f41247ec667e6',
        'EfpD+gHcLOn+ZxUGqJAVMBUmrXg4iUmqGRfq0tK78Zzl4unb5w==',
        (SELECT id FROM roles WHERE name = 'ADMIN'),
        'ACTIVE'
    ),
    (
        '00000000-0000-4000-8000-000000000002',
        'ymINtN76OVRfW1PvIBzi+J3dhBrp35nB0gUQ65PpFn4apYy/KpqvlR2eTQ==',
        'ULa0xtTtRHOWOFwi1CAws0rd398d9JhAQsO7m3FJoOxkFiTUlq+XpXPlp4bVc9sdroA=',
        '3c42f43ca98a926d5da41c3a69693f9b587139cbfa69371915b1870031f34311',
        'Zw+NH7AOHqWGEQpp8iLI6kVMWgQBwnY5EuI0N1RiJzuLY2eRQY/f',
        (SELECT id FROM roles WHERE name = 'CLIENTE'),
        'ACTIVE'
    )
ON CONFLICT (email_hash) DO NOTHING;

COMMIT;

