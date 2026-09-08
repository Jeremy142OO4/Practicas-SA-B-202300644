import json
import os

import aio_pika

from models.pago import ProcesarPagoRequest


class PaymentConsumer:

    def __init__(self, servicio, ordenes):
        self.servicio = servicio
        self.ordenes = ordenes
        self.url = os.getenv(
            "RABBITMQ_URL",
            "amqp://guest:guest@rabbitmq:5672/",
        )
        self.queue_name = os.getenv(
            "PAYMENT_QUEUE",
            "payment.requests",
        )
        self.connection = None
        self.channel = None
        self.consumer_tag = None

    async def conectar(self):
        self.connection = await aio_pika.connect_robust(
            self.url
        )
        self.channel = await self.connection.channel()
        await self.channel.set_qos(prefetch_count=1)
        cola = await self.channel.declare_queue(
            self.queue_name,
            durable=True,
        )
        self.consumer_tag = await cola.consume(
            self.procesar_mensaje
        )

    async def cerrar(self):
        if self.connection is not None:
            await self.connection.close()

    async def procesar_mensaje(self, mensaje):
        async with mensaje.process(requeue=True):
            evento = json.loads(
                mensaje.body.decode("utf-8")
            )
            datos = ProcesarPagoRequest(
                id_orden=evento["id_orden"],
                id_metodo_pago=(
                    evento["id_metodo_pago"]
                ),
                monto=evento["monto"],
            )
            resultado = await self.servicio.procesar_pago(
                datos
            )
            await self.ordenes.actualizar_estado_pago(
                datos.id_orden,
                resultado["pago"]["estado"],
            )
