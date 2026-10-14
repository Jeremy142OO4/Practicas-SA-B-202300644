package main

import (
	"context"
	"log"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"p4/api-gateway/config"
	"p4/api-gateway/handlers"
	"p4/api-gateway/middleware"
	"p4/api-gateway/routes"
)

func main() {

	_ = godotenv.Load()

	configuracion :=
		config.CargarConfiguracion()

	shutdown, err := iniciarTelemetria(context.Background(), "api-gateway")
	if err != nil {
		log.Fatal("No se pudo iniciar OpenTelemetry: ", err)
	}
	defer shutdown(context.Background())

	app := fiber.New(
		fiber.Config{
			AppName: "API Gateway",
		},
	)
	app.Use(otelfiber.Middleware())
	app.Use(registrarMetricas("api-gateway"))
	app.Get("/metrics", endpointMetricas())

	// CORS
	app.Use(
		middleware.ConfigurarCORS(),
	)

	// Health
	app.Get(
		"/health",
		func(ctx *fiber.Ctx) error {

			return ctx.JSON(
				fiber.Map{
					"servicio": "api-gateway",
					"estado":   "ok",
				},
			)
		},
	)

	// Proxy
	proxyHandler :=
		handlers.NuevoProxyHandler()

	// Rutas
	routes.RegistrarRutas(
		app,
		configuracion,
		proxyHandler,
	)

	log.Println(
		"API Gateway ejecutándose en puerto",
		configuracion.Puerto,
	)

	log.Fatal(
		app.Listen(
			":" + configuracion.Puerto,
		),
	)
}
