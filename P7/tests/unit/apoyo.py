from ..conftest import cargar_modulo


ordenes = cargar_modulo("servicio_ordenes", "P4/orders-ms/services/order_service.py")
pagos = cargar_modulo("servicio_pagos", "P4/payments-ms/services/payment_service.py")


class SimuladorAsincrono:
    def __init__(self, **resultados):
        self.resultados = resultados
        self.llamadas = []

    def __getattr__(self, nombre):
        async def metodo(*argumentos, **argumentos_nombrados):
            self.llamadas.append((nombre, argumentos, argumentos_nombrados))
            resultado = self.resultados.get(nombre)
            return resultado(*argumentos, **argumentos_nombrados) if callable(resultado) else resultado

        return metodo
