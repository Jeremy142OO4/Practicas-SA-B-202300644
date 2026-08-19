CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS catalog;


CREATE TABLE IF NOT EXISTS catalog.categorias (

    id_categoria SERIAL PRIMARY KEY,

    nombre VARCHAR(100) NOT NULL,

    descripcion VARCHAR(255)
);


CREATE TABLE IF NOT EXISTS catalog.videojuegos (

    id_videojuego UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    titulo VARCHAR(150) NOT NULL,

    descripcion VARCHAR(500),

    id_categoria INT NOT NULL,

    desarrollador VARCHAR(150) NOT NULL,

    fecha_lanzamiento DATE NOT NULL,

    estado VARCHAR(20) NOT NULL
        DEFAULT 'ACTIVO',

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_videojuego_categoria
        FOREIGN KEY (id_categoria)
        REFERENCES catalog.categorias(id_categoria)
);


CREATE TABLE IF NOT EXISTS catalog.plataformas (

    id_plataforma SERIAL PRIMARY KEY,

    nombre VARCHAR(100) NOT NULL,

    descripcion VARCHAR(255)
);


CREATE TABLE IF NOT EXISTS catalog.videojuego_plataforma (

    id_videojuego_plataforma UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_videojuego UUID NOT NULL,

    id_plataforma INT NOT NULL,

    precio DECIMAL(10,2) NOT NULL,

    stock INT NOT NULL,

    activo BOOLEAN NOT NULL
        DEFAULT TRUE,

    CONSTRAINT fk_vp_videojuego
        FOREIGN KEY (id_videojuego)
        REFERENCES catalog.videojuegos(id_videojuego),

    CONSTRAINT fk_vp_plataforma
        FOREIGN KEY (id_plataforma)
        REFERENCES catalog.plataformas(id_plataforma),

    CONSTRAINT unique_videojuego_plataforma
        UNIQUE (
            id_videojuego,
            id_plataforma
        ),

    CONSTRAINT check_stock
        CHECK (stock >= 0),

    CONSTRAINT check_precio
        CHECK (precio >= 0)
);