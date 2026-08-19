from contextlib import asynccontextmanager
from pathlib import Path
import sys

from fastapi import FastAPI
from strawberry.fastapi import GraphQLRouter

from clients.catalog_client import CatalogClient
from clients.payment_client import PaymentClient
from config.database import (
    cerrar_base_datos,
    conectar_base_datos,
)
from controllers.order_controller import (
    OrderController,
)
sys.path.insert(0, str(Path(__file__).parent / "graphql"))
from schema import crear_esquema
from repositories.order_repository import (
    OrderRepository,
)
from services.order_service import (
    OrderService,
)


repositorio = OrderRepository()

catalogo = CatalogClient()

pagos = PaymentClient()

servicio = OrderService(
    repositorio,
    catalogo,
    pagos,
)

controlador = OrderController(
    servicio
)

esquema = crear_esquema(
    controlador
)

graphql_app = GraphQLRouter(
    esquema
)


@asynccontextmanager
async def lifespan(app: FastAPI):

    pool = await conectar_base_datos()

    repositorio.establecer_pool(
        pool
    )

    print(
        "Orders-ms conectado a PostgreSQL"
    )

    yield

    await cerrar_base_datos(
        pool
    )


app = FastAPI(
    title="Orders-ms",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/health")
async def health():

    return {
        "servicio": "orders-ms",
        "estado": "ok",
    }


app.include_router(
    graphql_app,
    prefix="/graphql",
)
