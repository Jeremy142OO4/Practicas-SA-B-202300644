package controlador

import (
	"github.com/gofiber/fiber/v2"

	respuesta "practica2/backend/internal/shared/response"
)

type Protegidas struct{}

func NuevaProtegidas() *Protegidas { return &Protegidas{} }

func (controlador *Protegidas) Administrador(ctx *fiber.Ctx) error {
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Acceso exclusivo autorizado", fiber.Map{
		"route": "admin", "role": ctx.Locals("usuario_rol"),
	})
}

func (controlador *Protegidas) Compartida(ctx *fiber.Ctx) error {
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Acceso compartido autorizado", fiber.Map{
		"route": "shared", "role": ctx.Locals("usuario_rol"),
	})
}
