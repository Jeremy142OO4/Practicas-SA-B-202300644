package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"p4/auth-ms/models"
	"p4/auth-ms/repositories"
	"p4/auth-ms/services"
)

type AuthController struct {
	Servicio *services.AuthService
}

func NuevoAuthController(
	servicio *services.AuthService,
) *AuthController {

	return &AuthController{
		Servicio: servicio,
	}
}

func (c *AuthController) RegistrarUsuario(
	ctx *fiber.Ctx,
) error {

	var datos models.RegistroRequest

	if err := ctx.BodyParser(&datos); err != nil {

		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": "datos inválidos",
			},
		)
	}

	usuario, err :=
		c.Servicio.RegistrarUsuario(
			datos,
		)

	if err != nil {

		if errors.Is(
			err,
			repositories.ErrCorreoRegistrado,
		) {

			return ctx.Status(
				fiber.StatusConflict,
			).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}

		if errors.Is(
			err,
			services.ErrDatosInvalidos,
		) {

			return ctx.Status(
				fiber.StatusBadRequest,
			).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}

		return ctx.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"error": "error interno",
			},
		)
	}

	return ctx.Status(
		fiber.StatusCreated,
	).JSON(
		fiber.Map{
			"mensaje": "usuario registrado",
			"usuario": usuario,
		},
	)
}

func (c *AuthController) IniciarSesion(
	ctx *fiber.Ctx,
) error {

	var datos models.LoginRequest

	if err := ctx.BodyParser(&datos); err != nil {

		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": "datos inválidos",
			},
		)
	}

	respuesta, err :=
		c.Servicio.IniciarSesion(
			datos,
		)

	if err != nil {

		if errors.Is(
			err,
			services.ErrCredencialesInvalidas,
		) ||
			errors.Is(
				err,
				services.ErrUsuarioInactivo,
			) {

			return ctx.Status(
				fiber.StatusUnauthorized,
			).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}

		return ctx.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"error": "error interno",
			},
		)
	}

	return ctx.JSON(respuesta)
}
