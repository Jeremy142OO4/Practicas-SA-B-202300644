class PagoRepository:

    def __init__(self):
        self.pool = None

    def establecer_pool(self, pool):
        self.pool = pool

    async def listar_metodos_pago(self):
        consulta = """
            SELECT
                id_metodo_pago,
                nombre,
                descripcion,
                activo
            FROM payments.metodos_pago
            WHERE activo = TRUE
            ORDER BY id_metodo_pago
        """

        filas = await self.pool.fetch(consulta)

        return [dict(fila) for fila in filas]

    async def obtener_metodo_pago(self, id_metodo_pago: int):
        consulta = """
            SELECT
                id_metodo_pago,
                nombre,
                descripcion,
                activo
            FROM payments.metodos_pago
            WHERE id_metodo_pago = $1
        """

        fila = await self.pool.fetchrow(
            consulta,
            id_metodo_pago,
        )

        if fila is None:
            return None

        return dict(fila)

    async def guardar_pago_con_transaccion(
        self,
        id_orden,
        id_metodo_pago: int,
        monto,
        codigo_transaccion: str,
        estado: str,
        mensaje: str,
    ):

        async with self.pool.acquire() as conexion:

            async with conexion.transaction():

                consulta_pago = """
                    INSERT INTO payments.pagos (
                        id_metodo_pago,
                        id_orden,
                        monto,
                        estado
                    )
                    VALUES ($1, $2, $3, 'PENDIENTE')
                    RETURNING
                        id_pago,
                        id_metodo_pago,
                        id_orden,
                        monto,
                        estado,
                        created_at,
                        updated_at
                """

                pago = await conexion.fetchrow(
                    consulta_pago,
                    id_metodo_pago,
                    id_orden,
                    monto,
                )

                consulta_transaccion = """
                    INSERT INTO payments.transacciones_pago (
                        id_pago,
                        codigo_transaccion,
                        estado,
                        mensaje
                    )
                    VALUES ($1, $2, $3, $4)
                    RETURNING
                        id_transaccion,
                        id_pago,
                        codigo_transaccion,
                        estado,
                        mensaje,
                        fecha_transaccion
                """

                transaccion = await conexion.fetchrow(
                    consulta_transaccion,
                    pago["id_pago"],
                    codigo_transaccion,
                    estado,
                    mensaje,
                )

                consulta_actualizar = """
                    UPDATE payments.pagos
                    SET
                        estado = $2,
                        updated_at = CURRENT_TIMESTAMP
                    WHERE id_pago = $1
                    RETURNING
                        id_pago,
                        id_metodo_pago,
                        id_orden,
                        monto,
                        estado,
                        created_at,
                        updated_at
                """

                pago_actualizado = await conexion.fetchrow(
                    consulta_actualizar,
                    pago["id_pago"],
                    estado,
                )

        return {
            "pago": dict(pago_actualizado),
            "transaccion": dict(transaccion),
        }

    async def obtener_pago(self, id_pago):
        consulta = """
            SELECT
                id_pago,
                id_metodo_pago,
                id_orden,
                monto,
                estado,
                created_at,
                updated_at
            FROM payments.pagos
            WHERE id_pago = $1
        """

        fila = await self.pool.fetchrow(
            consulta,
            id_pago,
        )

        if fila is None:
            return None

        return dict(fila)

    async def listar_pagos_orden(self, id_orden):
        consulta = """
            SELECT
                id_pago,
                id_metodo_pago,
                id_orden,
                monto,
                estado,
                created_at,
                updated_at
            FROM payments.pagos
            WHERE id_orden = $1
            ORDER BY created_at DESC
        """

        filas = await self.pool.fetch(
            consulta,
            id_orden,
        )

        return [dict(fila) for fila in filas]