package main

import (
	"log"
	"os"

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

	controlador :=
		controllers.NuevoCatalogController(
			esquema,
		)

	app := fiber.New()

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
