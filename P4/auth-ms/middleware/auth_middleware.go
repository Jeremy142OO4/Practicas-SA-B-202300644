package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"p4/auth-ms/config"
)

func Autenticar(
	gestorJWT *config.GestorJWT,
) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		authorization :=
			ctx.Get("Authorization")

		if authorization == "" {

			return ctx.Status(
				fiber.StatusUnauthorized,
			).JSON(
				fiber.Map{
					"error": "token requerido",
				},
			)
		}

		partes := strings.Split(
			authorization,
			" ",
		)

		if len(partes) != 2 ||
			partes[0] != "Bearer" {

			return ctx.Status(
				fiber.StatusUnauthorized,
			).JSON(
				fiber.Map{
					"error": "token inválido",
				},
			)
		}

		claims, err :=
			gestorJWT.ValidarToken(
				partes[1],
			)

		if err != nil {

			return ctx.Status(
				fiber.StatusUnauthorized,
			).JSON(
				fiber.Map{
					"error": "token inválido",
				},
			)
		}

		ctx.Locals(
			"id_usuario",
			claims.IDUsuario,
		)

		ctx.Locals(
			"rol",
			claims.Rol,
		)

		return ctx.Next()
	}
}
