package controlador

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"practica2/backend/internal/modules/users/dto"
	repositorio "practica2/backend/internal/modules/users/repository"
	servicio "practica2/backend/internal/modules/users/service"
	respuesta "practica2/backend/internal/shared/response"
)

type Usuarios struct {
	consulta servicio.ConsultaUsuarios
	estados  servicio.CambioEstados
}

func NuevosUsuarios(consulta servicio.ConsultaUsuarios, estados servicio.CambioEstados) (*Usuarios, error) {
	if consulta == nil || estados == nil {
		return nil, errors.New("los servicios de usuarios son obligatorios")
	}
	return &Usuarios{consulta: consulta, estados: estados}, nil
}

func (controlador *Usuarios) Listar(ctx *fiber.Ctx) error {
	usuarios, err := controlador.consulta.Listar(ctx.Context())
	if err != nil {
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudieron listar los usuarios", "USER_LIST_ERROR")
	}
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Usuarios obtenidos", usuarios)
}

func (controlador *Usuarios) Buscar(ctx *fiber.Ctx) error {
	usuario, err := controlador.consulta.Buscar(ctx.Context(), ctx.Params("id"))
	switch {
	case err == nil:
		return respuesta.Exitosa(ctx, fiber.StatusOK, "Usuario obtenido", usuario)
	case errors.Is(err, servicio.ErrIDInvalido):
		return respuesta.Fallo(ctx, fiber.StatusBadRequest, err.Error(), "INVALID_USER_ID")
	case errors.Is(err, repositorio.ErrUsuarioNoExiste):
		return respuesta.Fallo(ctx, fiber.StatusNotFound, "Usuario no encontrado", "USER_NOT_FOUND")
	default:
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudo consultar el usuario", "USER_QUERY_ERROR")
	}
}

func (controlador *Usuarios) CambiarEstado(ctx *fiber.Ctx) error {
	var solicitud dto.CambioEstado
	if err := ctx.BodyParser(&solicitud); err != nil {
		return respuesta.Fallo(ctx, fiber.StatusBadRequest, "Solicitud invalida", "INVALID_BODY")
	}
	administradorID, _ := ctx.Locals("usuario_id").(string)
	usuario, err := controlador.estados.Cambiar(ctx.Context(), administradorID, ctx.Params("id"), solicitud.Estado)
	switch {
	case err == nil:
		return respuesta.Exitosa(ctx, fiber.StatusOK, "Estado actualizado", usuario)
	case errors.Is(err, servicio.ErrIDInvalido), errors.Is(err, servicio.ErrEstadoInvalido):
		return respuesta.Fallo(ctx, fiber.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
	case errors.Is(err, servicio.ErrCambioPropio):
		return respuesta.Fallo(ctx, fiber.StatusConflict, err.Error(), "SELF_STATUS_CHANGE")
	case errors.Is(err, repositorio.ErrUsuarioNoExiste):
		return respuesta.Fallo(ctx, fiber.StatusNotFound, "Usuario no encontrado", "USER_NOT_FOUND")
	default:
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudo actualizar el usuario", "USER_UPDATE_ERROR")
	}
}
