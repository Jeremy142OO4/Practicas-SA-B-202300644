package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"practica2/backend/internal/shared/cookies"
	token "practica2/backend/internal/shared/jwt"
	respuesta "practica2/backend/internal/shared/response"
)

func Autenticar(gestorCookies *cookies.GestorSesion, validador token.Validador) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		valor := gestorCookies.Leer(ctx)
		if valor == "" {
			return respuesta.Fallo(ctx, fiber.StatusUnauthorized, "Debes iniciar sesion", "SESSION_REQUIRED")
		}
		datos, err := validador.Validar(valor)
		if err != nil {
			codigo := "INVALID_SESSION"
			mensaje := "La sesion no es valida"
			if errors.Is(err, token.ErrTokenExpirado) {
				codigo = "SESSION_EXPIRED"
				mensaje = "La sesion expiro y puede renovarse"
			}
			return respuesta.Fallo(ctx, fiber.StatusUnauthorized, mensaje, codigo)
		}
		ctx.Locals("usuario_id", datos.ID)
		ctx.Locals("usuario_rol", datos.Rol)
		ctx.Locals("sesion_expira", datos.ExpiraEn)
		ctx.Locals("sesion_renovable_hasta", datos.RenovableHasta)
		return ctx.Next()
	}
}
