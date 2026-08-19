from uuid import UUID

import strawberry

from resolver import (
    CarritoTipo,
    OrdenTipo,
    convertir_carrito,
    convertir_orden,
)


def crear_esquema(controlador):

    @strawberry.type
    class Query:

        @strawberry.field
        async def carrito(
            self,
            id_usuario: strawberry.ID,
        ) -> CarritoTipo:

            datos = (
                await controlador
                .obtener_carrito(
                    UUID(str(id_usuario))
                )
            )

            return convertir_carrito(
                datos
            )

        @strawberry.field
        async def orden(
            self,
            id_orden: strawberry.ID,
        ) -> OrdenTipo:

            datos = (
                await controlador
                .obtener_orden(
                    UUID(str(id_orden))
                )
            )

            return convertir_orden(
                datos
            )

        @strawberry.field
        async def ordenes(
            self,
            id_usuario: strawberry.ID,
        ) -> list[OrdenTipo]:

            datos = (
                await controlador
                .listar_ordenes(
                    UUID(str(id_usuario))
                )
            )

            return [
                convertir_orden(orden)
                for orden in datos
            ]

    @strawberry.type
    class Mutation:

        @strawberry.mutation
        async def agregar_producto(
            self,
            id_usuario: strawberry.ID,
            id_videojuego_plataforma:
                strawberry.ID,
            cantidad: int,
        ) -> CarritoTipo:

            datos = (
                await controlador
                .agregar_producto(
                    UUID(str(id_usuario)),
                    UUID(
                        str(
                            id_videojuego_plataforma
                        )
                    ),
                    cantidad,
                )
            )

            return convertir_carrito(
                datos
            )

        @strawberry.mutation
        async def eliminar_producto(
            self,
            id_usuario: strawberry.ID,
            id_detalle_carrito:
                strawberry.ID,
        ) -> CarritoTipo:

            datos = (
                await controlador
                .eliminar_producto(
                    UUID(str(id_usuario)),
                    UUID(
                        str(
                            id_detalle_carrito
                        )
                    ),
                )
            )

            return convertir_carrito(
                datos
            )

        @strawberry.mutation
        async def crear_orden(
            self,
            id_carrito: strawberry.ID,
        ) -> OrdenTipo:

            datos = (
                await controlador
                .crear_orden(
                    UUID(str(id_carrito))
                )
            )

            return convertir_orden(
                datos
            )

        @strawberry.mutation
        async def pagar_orden(
            self,
            id_orden: strawberry.ID,
            id_metodo_pago: int,
        ) -> OrdenTipo:

            datos = (
                await controlador
                .pagar_orden(
                    UUID(str(id_orden)),
                    id_metodo_pago,
                )
            )

            return convertir_orden(
                datos
            )

    return strawberry.Schema(
        query=Query,
        mutation=Mutation,
    )
