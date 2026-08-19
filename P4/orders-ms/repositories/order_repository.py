class OrderRepository:

    def __init__(self):
        self.pool = None

    def establecer_pool(self, pool):
        self.pool = pool

    async def obtener_carrito_activo(
        self,
        id_usuario,
    ):

        consulta = """
            SELECT
                id_carrito,
                id_usuario,
                estado,
                created_at,
                updated_at
            FROM orders.carritos
            WHERE
                id_usuario = $1
                AND estado = 'ACTIVO'
            ORDER BY created_at DESC
            LIMIT 1
        """

        fila = await self.pool.fetchrow(
            consulta,
            id_usuario,
        )

        return (
            dict(fila)
            if fila
            else None
        )

    async def obtener_carrito(
        self,
        id_carrito,
    ):

        consulta = """
            SELECT
                id_carrito,
                id_usuario,
                estado,
                created_at,
                updated_at
            FROM orders.carritos
            WHERE id_carrito = $1
        """

        fila = await self.pool.fetchrow(
            consulta,
            id_carrito,
        )

        return (
            dict(fila)
            if fila
            else None
        )

    async def crear_carrito(
        self,
        id_usuario,
    ):

        consulta = """
            INSERT INTO orders.carritos (
                id_usuario,
                estado
            )
            VALUES ($1, 'ACTIVO')
            RETURNING
                id_carrito,
                id_usuario,
                estado,
                created_at,
                updated_at
        """

        fila = await self.pool.fetchrow(
            consulta,
            id_usuario,
        )

        return dict(fila)

    async def listar_detalles_carrito(
        self,
        id_carrito,
    ):

        consulta = """
            SELECT
                id_detalle_carrito,
                id_carrito,
                id_videojuego_plataforma,
                cantidad,
                precio_unitario
            FROM orders.detalle_carrito
            WHERE id_carrito = $1
            ORDER BY id_detalle_carrito
        """

        filas = await self.pool.fetch(
            consulta,
            id_carrito,
        )

        return [
            dict(fila)
            for fila in filas
        ]

    async def agregar_producto(
        self,
        id_carrito,
        id_videojuego_plataforma,
        cantidad,
        precio_unitario,
    ):

        consulta = """
            INSERT INTO orders.detalle_carrito (
                id_carrito,
                id_videojuego_plataforma,
                cantidad,
                precio_unitario
            )
            VALUES ($1, $2, $3, $4)

            ON CONFLICT (
                id_carrito,
                id_videojuego_plataforma
            )

            DO UPDATE SET
                cantidad =
                    detalle_carrito.cantidad
                    + EXCLUDED.cantidad,

                precio_unitario =
                    EXCLUDED.precio_unitario

            RETURNING
                id_detalle_carrito,
                id_carrito,
                id_videojuego_plataforma,
                cantidad,
                precio_unitario
        """

        fila = await self.pool.fetchrow(
            consulta,
            id_carrito,
            id_videojuego_plataforma,
            cantidad,
            precio_unitario,
        )

        return dict(fila)

    async def eliminar_producto(
        self,
        id_detalle_carrito,
    ):

        resultado = await self.pool.execute(
            """
            DELETE FROM orders.detalle_carrito
            WHERE id_detalle_carrito = $1
            """,
            id_detalle_carrito,
        )

        return resultado == "DELETE 1"

    async def crear_orden(
        self,
        id_carrito,
        id_usuario,
        total,
        detalles,
    ):

        async with self.pool.acquire() as conexion:

            async with conexion.transaction():

                orden = await conexion.fetchrow(
                    """
                    INSERT INTO orders.ordenes (
                        id_carrito,
                        id_usuario,
                        total,
                        estado
                    )
                    VALUES (
                        $1,
                        $2,
                        $3,
                        'PENDIENTE'
                    )
                    RETURNING
                        id_orden,
                        id_carrito,
                        id_usuario,
                        total,
                        estado,
                        created_at,
                        updated_at
                    """,
                    id_carrito,
                    id_usuario,
                    total,
                )

                for detalle in detalles:

                    await conexion.execute(
                        """
                        INSERT INTO
                        orders.detalle_orden (
                            id_orden,
                            id_videojuego_plataforma,
                            nombre_videojuego,
                            plataforma,
                            cantidad,
                            precio_unitario,
                            subtotal
                        )
                        VALUES (
                            $1,
                            $2,
                            $3,
                            $4,
                            $5,
                            $6,
                            $7
                        )
                        """,
                        orden["id_orden"],
                        detalle[
                            "id_videojuego_plataforma"
                        ],
                        detalle[
                            "nombre_videojuego"
                        ],
                        detalle["plataforma"],
                        detalle["cantidad"],
                        detalle[
                            "precio_unitario"
                        ],
                        detalle["subtotal"],
                    )

                await conexion.execute(
                    """
                    UPDATE orders.carritos
                    SET
                        estado = 'PROCESADO',
                        updated_at =
                            CURRENT_TIMESTAMP
                    WHERE id_carrito = $1
                    """,
                    id_carrito,
                )

        return dict(orden)

    async def obtener_orden(
        self,
        id_orden,
    ):

        fila = await self.pool.fetchrow(
            """
            SELECT
                id_orden,
                id_carrito,
                id_usuario,
                total,
                estado,
                created_at,
                updated_at
            FROM orders.ordenes
            WHERE id_orden = $1
            """,
            id_orden,
        )

        return (
            dict(fila)
            if fila
            else None
        )

    async def listar_detalles_orden(
        self,
        id_orden,
    ):

        filas = await self.pool.fetch(
            """
            SELECT
                id_detalle_orden,
                id_orden,
                id_videojuego_plataforma,
                nombre_videojuego,
                plataforma,
                cantidad,
                precio_unitario,
                subtotal
            FROM orders.detalle_orden
            WHERE id_orden = $1
            ORDER BY id_detalle_orden
            """,
            id_orden,
        )

        return [
            dict(fila)
            for fila in filas
        ]

    async def listar_ordenes(
        self,
        id_usuario,
    ):

        filas = await self.pool.fetch(
            """
            SELECT
                id_orden,
                id_carrito,
                id_usuario,
                total,
                estado,
                created_at,
                updated_at
            FROM orders.ordenes
            WHERE id_usuario = $1
            ORDER BY created_at DESC
            """,
            id_usuario,
        )

        return [
            dict(fila)
            for fila in filas
        ]

    async def actualizar_estado_orden(
        self,
        id_orden,
        estado,
    ):

        fila = await self.pool.fetchrow(
            """
            UPDATE orders.ordenes
            SET
                estado = $2,
                updated_at = CURRENT_TIMESTAMP
            WHERE id_orden = $1
            RETURNING
                id_orden,
                id_carrito,
                id_usuario,
                total,
                estado,
                created_at,
                updated_at
            """,
            id_orden,
            estado,
        )

        return (
            dict(fila)
            if fila
            else None
        )