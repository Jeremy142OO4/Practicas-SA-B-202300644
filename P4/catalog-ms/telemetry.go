package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_server_requests_total", Help: "Total de solicitudes HTTP."},
		[]string{"service", "method", "route", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_server_request_duration_seconds", Help: "Duracion de solicitudes HTTP."},
		[]string{"service", "method", "route"},
	)
)

func iniciarTelemetria(ctx context.Context, service string) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "otel-collector.observability.svc.cluster.local:4317"
	}
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	resource, err := sdkresource.New(ctx, sdkresource.WithAttributes(
		attribute.String("service.name", service),
		attribute.String("deployment.environment", os.Getenv("ENVIRONMENT")),
	))
	if err != nil {
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(resource))
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	prometheus.MustRegister(requestTotal, requestDuration)
	return provider.Shutdown, nil
}

func registrarMetricas(service string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()
		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}
		requestTotal.WithLabelValues(service, c.Method(), route, strconv.Itoa(status)).Inc()
		requestDuration.WithLabelValues(service, c.Method(), route).Observe(time.Since(start).Seconds())
		spanContext := trace.SpanContextFromContext(c.UserContext())
		log.Printf(`{"service":"%s","event":"http_request","method":"%s","route":"%s","status":%d,"duration_ms":%d,"trace_id":"%s"}`,
			service, c.Method(), route, status, time.Since(start).Milliseconds(), spanContext.TraceID().String())
		return err
	}
}

func endpointMetricas() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}
