package respuesta

import "github.com/gofiber/fiber/v2"

type Respuesta struct {
	Exito   bool        `json:"success"`
	Mensaje string      `json:"message"`
	Datos   interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Codigo string `json:"code"`
}

func Exitosa(ctx *fiber.Ctx, estado int, mensaje string, datos interface{}) error {
	return ctx.Status(estado).JSON(Respuesta{Exito: true, Mensaje: mensaje, Datos: datos})
}

func Fallo(ctx *fiber.Ctx, estado int, mensaje string, codigo string) error {
	return ctx.Status(estado).JSON(Respuesta{
		Exito:   false,
		Mensaje: mensaje,
		Error:   &Error{Codigo: codigo},
	})
}
