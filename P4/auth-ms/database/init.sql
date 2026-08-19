CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.roles (
    id_rol SERIAL PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL UNIQUE,
    descripcion VARCHAR(150)
);

CREATE TABLE IF NOT EXISTS auth.usuarios (
    id_usuario UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name_encrypted VARCHAR(500) NOT NULL,

    email_encrypted VARCHAR(500) NOT NULL,

    email_hash VARCHAR(64) NOT NULL UNIQUE,

    password_encrypted VARCHAR(500) NOT NULL,

    role_id INT NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVO',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_usuario_rol
        FOREIGN KEY (role_id)
        REFERENCES auth.roles(id_rol)
);

INSERT INTO auth.roles (
    id_rol,
    nombre,
    descripcion
)
VALUES
    (
        1,
        'ADMIN',
        'Administrador del sistema'
    ),
    (
        2,
        'CLIENTE',
        'Cliente de la tienda'
    )
ON CONFLICT (id_rol)
DO NOTHING;