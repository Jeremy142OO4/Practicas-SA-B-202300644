package middleware

import (
	"github.com/gofiber/fiber/v2"

	respuesta "practica2/backend/internal/shared/response"
)

func Autorizar(roles ...string) fiber.Handler {
	permitidos := make(map[string]struct{}, len(roles))
	for _, rol := range roles {
		permitidos[rol] = struct{}{}
	}
	return func(ctx *fiber.Ctx) error {
		rol, correcto := ctx.Locals("usuario_rol").(string)
		if !correcto {
			return respuesta.Fallo(ctx, fiber.StatusUnauthorized, "Debes iniciar sesion", "SESSION_REQUIRED")
		}
		if _, permitido := permitidos[rol]; !permitido {
			return respuesta.Fallo(ctx, fiber.StatusForbidden, "No tienes permiso para acceder", "FORBIDDEN")
		}
		return ctx.Next()
	}
}
