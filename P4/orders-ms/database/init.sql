CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS orders;


CREATE TABLE IF NOT EXISTS orders.carritos (

    id_carrito UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_usuario UUID NOT NULL,

    estado VARCHAR(20) NOT NULL
        DEFAULT 'ACTIVO',

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS orders.detalle_carrito (

    id_detalle_carrito UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_carrito UUID NOT NULL,

    id_videojuego_plataforma UUID NOT NULL,

    cantidad INT NOT NULL,

    precio_unitario DECIMAL(10,2)
        NOT NULL,

    CONSTRAINT fk_detalle_carrito
        FOREIGN KEY (id_carrito)
        REFERENCES orders.carritos(
            id_carrito
        )
        ON DELETE CASCADE,

    CONSTRAINT unique_producto_carrito
        UNIQUE (
            id_carrito,
            id_videojuego_plataforma
        ),

    CONSTRAINT check_cantidad_carrito
        CHECK (cantidad > 0),

    CONSTRAINT check_precio_carrito
        CHECK (precio_unitario >= 0)
);


CREATE TABLE IF NOT EXISTS orders.ordenes (

    id_orden UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_carrito UUID NOT NULL UNIQUE,

    id_usuario UUID NOT NULL,

    total DECIMAL(10,2) NOT NULL,

    estado VARCHAR(20) NOT NULL
        DEFAULT 'PENDIENTE',

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_orden_carrito
        FOREIGN KEY (id_carrito)
        REFERENCES orders.carritos(
            id_carrito
        ),

    CONSTRAINT check_total
        CHECK (total >= 0)
);


CREATE TABLE IF NOT EXISTS orders.detalle_orden (

    id_detalle_orden UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_orden UUID NOT NULL,

    id_videojuego_plataforma UUID NOT NULL,

    nombre_videojuego VARCHAR(150)
        NOT NULL,

    plataforma VARCHAR(100)
        NOT NULL,

    cantidad INT NOT NULL,

    precio_unitario DECIMAL(10,2)
        NOT NULL,

    subtotal DECIMAL(10,2)
        NOT NULL,

    CONSTRAINT fk_detalle_orden
        FOREIGN KEY (id_orden)
        REFERENCES orders.ordenes(
            id_orden
        )
        ON DELETE CASCADE,

    CONSTRAINT check_cantidad_orden
        CHECK (cantidad > 0),

    CONSTRAINT check_precio_orden
        CHECK (precio_unitario >= 0),

    CONSTRAINT check_subtotal
        CHECK (subtotal >= 0)
);