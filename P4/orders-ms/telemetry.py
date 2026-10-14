import json
import logging
import os
import time

from fastapi import FastAPI, Request
from opentelemetry import propagate, trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.httpx import HTTPXClientInstrumentor
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from prometheus_client import Counter, Histogram, make_asgi_app


REQUESTS = Counter(
    "http_server_requests_total",
    "Total de solicitudes HTTP.",
    ["service", "method", "route", "status"],
)
REQUEST_DURATION = Histogram(
    "http_server_request_duration_seconds",
    "Duracion de solicitudes HTTP.",
    ["service", "method", "route"],
)
BUSINESS_EVENTS = Counter(
    "business_events_total",
    "Eventos de negocio procesados.",
    ["service", "event", "status"],
)
PAYMENTS_APPROVED = Counter("payments_approved_total", "Resultados de pago aprobado procesados.")
PAYMENTS_REJECTED = Counter("payments_rejected_total", "Resultados de pago rechazado procesados.")


def configure_telemetry(app: FastAPI, service: str) -> TracerProvider:
    logging.basicConfig(level=logging.INFO, format="%(message)s")
    logging.getLogger(service).setLevel(logging.INFO)
    endpoint = os.getenv(
        "OTEL_EXPORTER_OTLP_ENDPOINT",
        "otel-collector.observability.svc.cluster.local:4317",
    )
    provider = TracerProvider(
        resource=Resource.create(
            {
                "service.name": service,
                "deployment.environment": os.getenv("ENVIRONMENT", "dev"),
            }
        )
    )
    exporter = OTLPSpanExporter(endpoint=endpoint, insecure=True)
    provider.add_span_processor(BatchSpanProcessor(exporter))
    trace.set_tracer_provider(provider)
    propagate.set_global_textmap(TraceContextTextMapPropagator())
    FastAPIInstrumentor.instrument_app(app, tracer_provider=provider)
    HTTPXClientInstrumentor().instrument()

    @app.middleware("http")
    async def metrics_and_json_log(request: Request, call_next):
        started = time.perf_counter()
        response = await call_next(request)
        elapsed = time.perf_counter() - started
        route = request.scope.get("route")
        route_name = getattr(route, "path", request.url.path)
        status = str(response.status_code)
        REQUESTS.labels(service, request.method, route_name, status).inc()
        REQUEST_DURATION.labels(service, request.method, route_name).observe(elapsed)
        span_context = trace.get_current_span().get_span_context()
        logging.getLogger(service).info(
            json.dumps(
                {
                    "service": service,
                    "event": "http_request",
                    "method": request.method,
                    "route": route_name,
                    "status": response.status_code,
                    "duration_ms": round(elapsed * 1000, 2),
                    "trace_id": format(span_context.trace_id, "032x"),
                }
            )
        )
        return response

    app.mount("/metrics", make_asgi_app())
    return provider


def record_business_event(service: str, event: str, status: str) -> None:
    BUSINESS_EVENTS.labels(service, event, status).inc()
    if event == "payment_result" and status == "aprobado":
        PAYMENTS_APPROVED.inc()
    elif event == "payment_result" and status == "rechazado":
        PAYMENTS_REJECTED.inc()
