package controller

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/tu-usuario/academia-api/internal/modules/solicitudes/dto"
	"github.com/tu-usuario/academia-api/internal/modules/solicitudes/service"
	errorsapp "github.com/tu-usuario/academia-api/internal/shared/errors"
	"github.com/tu-usuario/academia-api/internal/shared/response"
)

type Controller struct {
	service service.Service
}

func NewController(
	service service.Service,
) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller *Controller) ObtenerTodas(
	c fiber.Ctx,
) error {
	solicitudes, err :=
		controller.service.ObtenerTodas(
			c.Context(),
		)

	if err != nil {
		return responderError(c, err)
	}

	return response.JSON(
		c,
		fiber.StatusOK,
		solicitudes,
	)
}

func (controller *Controller) Crear(
	c fiber.Ctx,
) error {
	var input dto.CrearSolicitudRequest

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			"INVALID_JSON",
			"el cuerpo JSON no es válido",
		)
	}

	solicitud, err :=
		controller.service.Crear(
			c.Context(),
			input,
		)

	if err != nil {
		return responderError(c, err)
	}

	return response.JSON(
		c,
		fiber.StatusCreated,
		solicitud,
	)
}

func (controller *Controller) Actualizar(
	c fiber.Ctx,
) error {
	var input dto.ActualizarSolicitudRequest

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			"INVALID_JSON",
			"el cuerpo JSON no es válido",
		)
	}

	solicitud, err :=
		controller.service.Actualizar(
			c.Context(),
			c.Params("id"),
			input,
		)

	if err != nil {
		return responderError(c, err)
	}

	return response.JSON(
		c,
		fiber.StatusOK,
		solicitud,
	)
}

func (controller *Controller) Eliminar(
	c fiber.Ctx,
) error {
	err := controller.service.Eliminar(
		c.Context(),
		c.Params("id"),
	)

	if err != nil {
		return responderError(c, err)
	}

	return c.SendStatus(
		fiber.StatusNoContent,
	)
}

func (controller *Controller) ActualizarEstado(
	c fiber.Ctx,
) error {
	var input dto.ActualizarEstadoRequest

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			"INVALID_JSON",
			"el cuerpo JSON no es válido",
		)
	}

	solicitud, err :=
		controller.service.ActualizarEstado(
			c.Context(),
			c.Params("id"),
			input,
		)

	if err != nil {
		return responderError(c, err)
	}

	return response.JSON(
		c,
		fiber.StatusOK,
		solicitud,
	)
}

func responderError(
	c fiber.Ctx,
	err error,
) error {
	var appErr *errorsapp.AppError

	if !errorsapp.Es(err, &appErr) {
		log.Printf(
			"error no controlado: %v",
			err,
		)

		return response.Error(
			c,
			fiber.StatusInternalServerError,
			string(errorsapp.TipoInterno),
			"error interno del servidor",
		)
	}

	switch appErr.Tipo {
	case errorsapp.TipoValidacion:
		return response.Error(
			c,
			fiber.StatusBadRequest,
			string(appErr.Tipo),
			appErr.Mensaje,
		)

	case errorsapp.TipoNoEncontrado:
		return response.Error(
			c,
			fiber.StatusNotFound,
			string(appErr.Tipo),
			appErr.Mensaje,
		)

	case errorsapp.TipoConflicto:
		return response.Error(
			c,
			fiber.StatusConflict,
			string(appErr.Tipo),
			appErr.Mensaje,
		)

	default:
		log.Printf(
			"error interno: %v",
			appErr.Causa,
		)

		return response.Error(
			c,
			fiber.StatusInternalServerError,
			string(errorsapp.TipoInterno),
			"error interno del servidor",
		)
	}
}
