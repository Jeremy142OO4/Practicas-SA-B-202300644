from contextlib import asynccontextmanager
from pathlib import Path
import sys

from fastapi import FastAPI
from dotenv import load_dotenv
from strawberry.fastapi import GraphQLRouter

load_dotenv()

from clients.catalog_client import CatalogClient
from clients.payment_publisher import PaymentPublisher
from config.database import (
    cerrar_base_datos,
    conectar_base_datos,
)
from controllers.order_controller import (
    OrderController,
)
from models.resultado_pago import ResultadoPagoRequest
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

pagos = PaymentPublisher()

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

    await pagos.conectar()

    print(
        "Orders-ms conectado a PostgreSQL"
    )

    yield

    await pagos.cerrar()

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


@app.post("/pagos/resultado")
async def actualizar_resultado_pago(
    resultado: ResultadoPagoRequest,
):
    return await servicio.actualizar_resultado_pago(
        resultado.model_dump()
    )
