package routes

import (
	"github.com/gofiber/fiber/v2"

	"p4/auth-ms/config"
	"p4/auth-ms/controllers"
	"p4/auth-ms/middleware"
)

func RegistrarRutasAuth(
	app *fiber.App,
	controlador *controllers.AuthController,
	gestorJWT *config.GestorJWT,
) {

	auth := app.Group("/auth")

	auth.Post(
		"/register",
		controlador.RegistrarUsuario,
	)

	auth.Post(
		"/login",
		controlador.IniciarSesion,
	)

	auth.Get(
		"/validate",
		middleware.Autenticar(gestorJWT),
		func(ctx *fiber.Ctx) error {

			return ctx.JSON(
				fiber.Map{
					"valido": true,

					"id_usuario": ctx.Locals(
						"id_usuario",
					),

					"rol": ctx.Locals(
						"rol",
					),
				},
			)
		},
	)
}
