import os

import httpx


class CatalogClient:

    def __init__(self):
        self.url = os.getenv(
            "CATALOG_URL",
            "http://catalog-ms:4002/graphql",
        )

    async def ejecutar(self, consulta, variables):
        async with httpx.AsyncClient(
            timeout=10.0
        ) as cliente:

            respuesta = await cliente.post(
                self.url,
                json={
                    "query": consulta,
                    "variables": variables,
                },
            )

            respuesta.raise_for_status()

            datos = respuesta.json()

            if datos.get("errors"):
                raise Exception(
                    datos["errors"][0]["message"]
                )

            return datos["data"]

    async def obtener_producto(
        self,
        id_videojuego_plataforma,
    ):

        consulta_presentacion = """
            query ObtenerPresentacion($id: ID!) {
                videojuegoPlataforma(id: $id) {
                    idVideojuegoPlataforma
                    idVideojuego
                    idPlataforma
                    precio
                    stock
                    activo
                    plataforma {
                        nombre
                    }
                }
            }
        """

        datos = await self.ejecutar(
            consulta_presentacion,
            {
                "id": str(
                    id_videojuego_plataforma
                )
            },
        )

        presentacion = datos[
            "videojuegoPlataforma"
        ]

        if presentacion is None:
            raise Exception(
                "producto no encontrado"
            )

        consulta_videojuego = """
            query ObtenerVideojuego($id: ID!) {
                videojuego(id: $id) {
                    idVideojuego
                    titulo
                }
            }
        """

        datos_videojuego = await self.ejecutar(
            consulta_videojuego,
            {
                "id": presentacion[
                    "idVideojuego"
                ]
            },
        )

        videojuego = datos_videojuego[
            "videojuego"
        ]

        return {
            "id_videojuego_plataforma":
                presentacion[
                    "idVideojuegoPlataforma"
                ],

            "id_videojuego":
                presentacion["idVideojuego"],

            "titulo":
                videojuego["titulo"],

            "plataforma":
                presentacion[
                    "plataforma"
                ]["nombre"],

            "precio":
                presentacion["precio"],

            "stock":
                presentacion["stock"],

            "activo":
                presentacion["activo"],
        }

    async def actualizar_stock(
        self,
        id_videojuego_plataforma,
        stock,
    ):

        consulta = """
            mutation ActualizarStock(
                $id: ID!,
                $stock: Int!
            ) {
                actualizarStock(
                    id: $id,
                    stock: $stock
                ) {
                    idVideojuegoPlataforma
                    stock
                }
            }
        """

        datos = await self.ejecutar(
            consulta,
            {
                "id": str(
                    id_videojuego_plataforma
                ),
                "stock": stock,
            },
        )

        return datos["actualizarStock"]