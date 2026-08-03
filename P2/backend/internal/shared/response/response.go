package response

import "github.com/gofiber/fiber/v2"

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code string `json:"code"`
}

func Success(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	return ctx.Status(status).JSON(Envelope{Success: true, Message: message, Data: data})
}

func Failure(ctx *fiber.Ctx, status int, message string, code string) error {
	return ctx.Status(status).JSON(Envelope{
		Success: false,
		Message: message,
		Error:   &ErrorBody{Code: code},
	})
}
