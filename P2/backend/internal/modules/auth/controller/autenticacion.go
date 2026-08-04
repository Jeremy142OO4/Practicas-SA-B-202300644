package controlador

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"practica2/backend/internal/modules/auth/dto"
	servicio "practica2/backend/internal/modules/auth/service"
	repositorio "practica2/backend/internal/modules/users/repository"
	"practica2/backend/internal/shared/cookies"
	token "practica2/backend/internal/shared/jwt"
	"practica2/backend/internal/shared/response"
	validacion "practica2/backend/internal/shared/validator"
)

type Autenticacion struct {
	registro servicio.RegistroUsuarios
	login    servicio.InicioSesiones
	cookies  *cookies.GestorSesion
	tokens   *token.Gestor
}

func NuevaAutenticacion(
	registro servicio.RegistroUsuarios,
	login servicio.InicioSesiones,
	gestorCookies *cookies.GestorSesion,
	gestorTokens *token.Gestor,
) (*Autenticacion, error) {
	if registro == nil || login == nil || gestorCookies == nil || gestorTokens == nil {
		return nil, errors.New("los servicios de autenticacion son obligatorios")
	}
	return &Autenticacion{registro: registro, login: login, cookies: gestorCookies, tokens: gestorTokens}, nil
}

func (controlador *Autenticacion) Renovar(ctx *fiber.Ctx) error {
	valor := controlador.cookies.Leer(ctx)
	if valor == "" {
		return respuesta.Fallo(ctx, fiber.StatusUnauthorized, "Sesion no encontrada", "SESSION_REQUIRED")
	}
	datos, err := controlador.tokens.ValidarRenovacion(valor)
	if err != nil {
		controlador.cookies.Eliminar(ctx)
		return respuesta.Fallo(ctx, fiber.StatusUnauthorized, "La sesion no se puede renovar", "SESSION_NOT_RENEWABLE")
	}
	nuevo, err := controlador.tokens.Generar(datos.ID, datos.Rol)
	if err != nil {
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudo renovar la sesion", "RENEWAL_ERROR")
	}
	controlador.cookies.Guardar(ctx, nuevo)
	nuevaSesion, err := controlador.tokens.Validar(nuevo)
	if err != nil {
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudo leer la sesion renovada", "RENEWAL_ERROR")
	}
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Sesion renovada", fiber.Map{
		"id": nuevaSesion.ID, "role": nuevaSesion.Rol,
		"expiresAt": nuevaSesion.ExpiraEn, "renewableUntil": nuevaSesion.RenovableHasta,
	})
}

func (controlador *Autenticacion) CerrarSesion(ctx *fiber.Ctx) error {
	controlador.cookies.Eliminar(ctx)
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Sesion cerrada", nil)
}

func (controlador *Autenticacion) VerSesion(ctx *fiber.Ctx) error {
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Sesion activa", fiber.Map{
		"id": ctx.Locals("usuario_id"), "role": ctx.Locals("usuario_rol"),
		"expiresAt":      ctx.Locals("sesion_expira"),
		"renewableUntil": ctx.Locals("sesion_renovable_hasta"),
	})
}

func (controlador *Autenticacion) VerAdministracion(ctx *fiber.Ctx) error {
	return respuesta.Exitosa(ctx, fiber.StatusOK, "Acceso de administrador autorizado", nil)
}

func (controlador *Autenticacion) Registrar(ctx *fiber.Ctx) error {
	var solicitud dto.SolicitudRegistro
	if err := ctx.BodyParser(&solicitud); err != nil {
		return respuesta.Fallo(ctx, fiber.StatusBadRequest, "Solicitud invalida", "INVALID_BODY")
	}

	usuario, err := controlador.registro.Registrar(ctx.Context(), solicitud)
	switch {
	case err == nil:
		return respuesta.Exitosa(ctx, fiber.StatusCreated, "Usuario registrado", usuario)
	case errors.Is(err, repositorio.ErrCorreoDuplicado):
		return respuesta.Fallo(ctx, fiber.StatusConflict, "El correo ya esta registrado", "EMAIL_ALREADY_EXISTS")
	case validacion.EsErrorRegistro(err):
		return respuesta.Fallo(ctx, fiber.StatusUnprocessableEntity, err.Error(), "VALIDATION_ERROR")
	default:
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudo registrar el usuario", "REGISTRATION_ERROR")
	}
}

func (controlador *Autenticacion) IniciarSesion(ctx *fiber.Ctx) error {
	var solicitud dto.SolicitudLogin
	if err := ctx.BodyParser(&solicitud); err != nil {
		return respuesta.Fallo(ctx, fiber.StatusBadRequest, "Solicitud invalida", "INVALID_BODY")
	}

	resultado, err := controlador.login.Iniciar(ctx.Context(), solicitud)
	switch {
	case err == nil:
		controlador.cookies.Guardar(ctx, resultado.Token)
		return respuesta.Exitosa(ctx, fiber.StatusOK, "Inicio de sesion exitoso", resultado.Usuario)
	case errors.Is(err, servicio.ErrCredenciales):
		return respuesta.Fallo(ctx, fiber.StatusUnauthorized, "Credenciales incorrectas", "INVALID_CREDENTIALS")
	case errors.Is(err, servicio.ErrUsuarioInactivo):
		return respuesta.Fallo(ctx, fiber.StatusForbidden, "Usuario inactivo", "USER_INACTIVE")
	default:
		return respuesta.Fallo(ctx, fiber.StatusInternalServerError, "No se pudo iniciar sesion", "LOGIN_ERROR")
	}
}
