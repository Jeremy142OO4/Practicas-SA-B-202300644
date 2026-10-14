package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func TestDemoraSolamenteEnRutaDeNegocio(t *testing.T) {
	t.Setenv("P10_DEMO_DELAY_MS", "40")
	app := fiber.New()
	app.Use(demoraDemostracion())
	app.All("/*", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	start := time.Now()
	response, err := app.Test(httptest.NewRequest("POST", "/orders/graphql", nil))
	if err != nil { t.Fatal(err) }
	response.Body.Close()
	if time.Since(start) < 40*time.Millisecond { t.Fatal("No se aplico la demora") }
	start = time.Now()
	response, err = app.Test(httptest.NewRequest("GET", "/health", nil))
	if err != nil { t.Fatal(err) }
	response.Body.Close()
	if time.Since(start) >= 40*time.Millisecond { t.Fatal("La demora afecto /health") }
}

func TestEtiquetasNoCambianConReutilizacionDeFiber(t *testing.T) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(requestTotal, requestDuration)
	app := fiber.New()
	app.Use(registrarMetricas("api-gateway"))
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	app.Post("/orders/graphql", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	for i := 0; i < 20; i++ {
		for _, target := range []struct{method, path string}{{"GET", "/health"}, {"POST", "/orders/graphql"}} {
			response, err := app.Test(httptest.NewRequest(target.method, target.path, nil))
			if err != nil { t.Fatal(err) }
			response.Body.Close()
		}
	}
	if _, err := registry.Gather(); err != nil { t.Fatal(err) }
}
