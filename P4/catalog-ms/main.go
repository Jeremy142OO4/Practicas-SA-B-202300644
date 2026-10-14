package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"p4/catalog-ms/config"
	"p4/catalog-ms/controllers"
	graphqlapi "p4/catalog-ms/graphql"
	"p4/catalog-ms/repositories"
	"p4/catalog-ms/services"
)

func main() {

	_ = godotenv.Load()

	bd, err := config.ConectarBaseDatos()

	if err != nil {

		log.Fatal(
			"Error al conectar con PostgreSQL: ",
			err,
		)
	}

	defer bd.Close()

	repositorio :=
		repositories.NuevoCatalogRepository(
			bd,
		)

	servicio :=
		services.NuevoCatalogService(
			repositorio,
		)

	esquema, err :=
		graphqlapi.CrearEsquema(
			servicio,
		)

	if err != nil {

		log.Fatal(
			"Error al crear esquema GraphQL: ",
			err,
		)
	}

	shutdown, err := iniciarTelemetria(context.Background(), "catalog-ms")
	if err != nil {
		log.Fatal("No se pudo iniciar OpenTelemetry: ", err)
	}
	defer shutdown(context.Background())

	controlador :=
		controllers.NuevoCatalogController(
			esquema,
		)

	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Use(registrarMetricas("catalog-ms"))
	app.Get("/metrics", endpointMetricas())

	app.Get(
		"/health",
		func(ctx *fiber.Ctx) error {

			return ctx.JSON(
				fiber.Map{
					"servicio": "catalog-ms",
					"estado":   "ok",
				},
			)
		},
	)

	app.Post(
		"/graphql",
		controlador.EjecutarGraphQL,
	)

	puerto := os.Getenv("PORT")

	if puerto == "" {
		puerto = "4002"
	}

	log.Println(
		"Catalog-ms ejecutándose en puerto",
		puerto,
	)

	log.Fatal(
		app.Listen(":" + puerto),
	)
}
