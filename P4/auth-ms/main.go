package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"p4/auth-ms/config"
	"p4/auth-ms/controllers"
	"p4/auth-ms/repositories"
	"p4/auth-ms/routes"
	"p4/auth-ms/services"
)

func main() {

	// Cargar variables de entorno
	_ = godotenv.Load()

	// Conectar base de datos
	bd, err := config.ConectarBaseDatos()

	if err != nil {
		log.Fatal(
			"Error al conectar con PostgreSQL: ",
			err,
		)
	}

	defer bd.Close()

	// JWT
	claveJWT := os.Getenv(
		"JWT_SECRET",
	)

	gestorJWT :=
		config.NuevoGestorJWT(
			claveJWT,
		)

	// Repository
	repositorio :=
		repositories.NuevoUsuarioRepository(
			bd,
		)

	// Service
	servicio :=
		services.NuevoAuthService(
			repositorio,
			gestorJWT,
			os.Getenv("ENCRYPTION_KEY"),
		)

	// Controller
	controlador :=
		controllers.NuevoAuthController(
			servicio,
		)

	// Fiber
	app := fiber.New()

	app.Get(
		"/health",
		func(ctx *fiber.Ctx) error {

			return ctx.JSON(
				fiber.Map{
					"servicio": "auth-ms",
					"estado":   "ok",
				},
			)
		},
	)

	// Rutas
	routes.RegistrarRutasAuth(
		app,
		controlador,
		gestorJWT,
	)

	puerto := os.Getenv("PORT")

	if puerto == "" {
		puerto = "4001"
	}

	log.Println(
		"Auth-ms ejecutándose en puerto",
		puerto,
	)

	log.Fatal(
		app.Listen(":" + puerto),
	)
}
