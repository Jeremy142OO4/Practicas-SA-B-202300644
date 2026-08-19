package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type ProxyHandler struct{}

func NuevoProxyHandler() *ProxyHandler {

	return &ProxyHandler{}
}

func (h *ProxyHandler) Reenviar(
	ctx *fiber.Ctx,
	destino string,
) error {

	query := string(
		ctx.Context().QueryArgs().QueryString(),
	)

	if query != "" {
		destino += "?" + query
	}

	err := proxy.Do(
		ctx,
		destino,
	)

	if err != nil {

		return ctx.Status(
			fiber.StatusBadGateway,
		).JSON(
			fiber.Map{
				"error": "microservicio no disponible",
			},
		)
	}

	ctx.Response().Header.Del(
		fiber.HeaderServer,
	)

	return nil
}
