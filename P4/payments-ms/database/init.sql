CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS payments;


CREATE TABLE IF NOT EXISTS payments.metodos_pago (

    id_metodo_pago SERIAL PRIMARY KEY,

    nombre VARCHAR(100) NOT NULL UNIQUE,

    descripcion VARCHAR(255),

    activo BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE TABLE IF NOT EXISTS payments.pagos (

    id_pago UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_metodo_pago INT NOT NULL,

    id_orden UUID NOT NULL,

    monto DECIMAL(10,2) NOT NULL,

    estado VARCHAR(30) NOT NULL
        DEFAULT 'PENDIENTE',

    created_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_pago_metodo
        FOREIGN KEY (id_metodo_pago)
        REFERENCES payments.metodos_pago(
            id_metodo_pago
        ),

    CONSTRAINT check_monto
        CHECK (monto > 0)
);


CREATE TABLE IF NOT EXISTS payments.transacciones_pago (

    id_transaccion UUID PRIMARY KEY
        DEFAULT gen_random_uuid(),

    id_pago UUID NOT NULL,

    codigo_transaccion VARCHAR(100)
        NOT NULL,

    estado VARCHAR(30)
        NOT NULL,

    mensaje VARCHAR(255),

    fecha_transaccion TIMESTAMP NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_transaccion_pago
        FOREIGN KEY (id_pago)
        REFERENCES payments.pagos(id_pago)
);


INSERT INTO payments.metodos_pago (
    nombre,
    descripcion
)
VALUES
    (
        'TARJETA',
        'Pago simulado con tarjeta'
    ),
    (
        'TRANSFERENCIA',
        'Pago simulado por transferencia'
    )
ON CONFLICT (nombre)
DO NOTHING;