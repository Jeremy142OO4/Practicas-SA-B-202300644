import json
import os
from uuid import uuid4

import aio_pika


class PaymentPublisher:

    def __init__(self):
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

    async def conectar(self):
        self.connection = await aio_pika.connect_robust(
            self.url
        )
        self.channel = await self.connection.channel()
        await self.channel.declare_queue(
            self.queue_name,
            durable=True,
        )

    async def cerrar(self):
        if self.connection is not None:
            await self.connection.close()

    async def publicar_pago(
        self,
        id_orden,
        id_metodo_pago,
        monto,
    ):
        evento = {
            "event_id": str(uuid4()),
            "id_orden": str(id_orden),
            "id_metodo_pago": id_metodo_pago,
            "monto": str(monto),
        }

        mensaje = aio_pika.Message(
            body=json.dumps(evento).encode("utf-8"),
            content_type="application/json",
            delivery_mode=(
                aio_pika.DeliveryMode.PERSISTENT
            ),
            message_id=evento["event_id"],
        )

        await self.channel.default_exchange.publish(
            mensaje,
            routing_key=self.queue_name,
        )

        return evento
