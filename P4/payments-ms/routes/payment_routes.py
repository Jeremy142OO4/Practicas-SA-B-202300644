from fastapi import APIRouter


def crear_rutas_pagos(controlador):

    router = APIRouter(
        prefix="/pagos",
        tags=["Pagos"],
    )

    router.add_api_route(
        "/metodos",
        controlador.listar_metodos_pago,
        methods=["GET"],
    )

    router.add_api_route(
        "/procesar",
        controlador.procesar_pago,
        methods=["POST"],
    )

    router.add_api_route(
        "/{id_pago}",
        controlador.obtener_pago,
        methods=["GET"],
    )

    router.add_api_route(
        "/orden/{id_orden}",
        controlador.listar_pagos_orden,
        methods=["GET"],
    )

    return router