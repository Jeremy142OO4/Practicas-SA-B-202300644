import os
from contextlib import asynccontextmanager

from fastapi import FastAPI
from dotenv import load_dotenv

load_dotenv()

from clients.order_client import OrderClient
from config.database import (
    cerrar_base_datos,
    conectar_base_datos,
)
from controllers.payment_controller import PagoController
from consumers.payment_consumer import PaymentConsumer
from repositories.payment_repository import PagoRepository
from routes.payment_routes import crear_rutas_pagos
from services.payment_service import PagoService


repositorio = PagoRepository()

servicio = PagoService(
    repositorio
)

ordenes = OrderClient()

consumidor = PaymentConsumer(
    servicio,
    ordenes,
)

controlador = PagoController(
    servicio
)


@asynccontextmanager
async def lifespan(app: FastAPI):

    pool = await conectar_base_datos()

    repositorio.establecer_pool(
        pool
    )

    await consumidor.conectar()

    print(
        "Payments-ms conectado a PostgreSQL"
    )

    yield

    await consumidor.cerrar()

    await cerrar_base_datos(
        pool
    )


app = FastAPI(
    title="Payments-ms",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/health")
async def health():

    return {
        "servicio": "payments-ms",
        "estado": "ok",
    }


app.include_router(
    crear_rutas_pagos(
        controlador
    )
)
