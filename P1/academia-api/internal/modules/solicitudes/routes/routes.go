package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/controller"
)

func Registrar(
	router fiber.Router,
	controller *controller.Controller,
) {
	solicitudes := router.Group(
		"/solicitudes",
	)

	solicitudes.Get(
		"/",
		controller.ObtenerTodas,
	)

	solicitudes.Post(
		"/",
		controller.Crear,
	)

	solicitudes.Put(
		"/:id",
		controller.Actualizar,
	)

	solicitudes.Delete(
		"/:id",
		controller.Eliminar,
	)

	solicitudes.Patch(
		"/:id/estado",
		controller.ActualizarEstado,
	)
}
