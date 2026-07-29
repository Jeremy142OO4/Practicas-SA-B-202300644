package response

import "github.com/gofiber/fiber/v3"

type exitosa struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type detalleError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type fallida struct {
	Success bool         `json:"success"`
	Error   detalleError `json:"error"`
}

func JSON(
	c fiber.Ctx,
	status int,
	data any,
) error {
	return c.Status(status).JSON(
		exitosa{
			Success: true,
			Data:    data,
		},
	)
}

func Error(
	c fiber.Ctx,
	status int,
	code string,
	message string,
) error {
	return c.Status(status).JSON(
		fallida{
			Success: false,
			Error: detalleError{
				Code:    code,
				Message: message,
			},
		},
	)
}
