package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ProxyHandler struct{}

func NuevoProxyHandler() *ProxyHandler {

	return &ProxyHandler{}
}

func (h *ProxyHandler) Reenviar(
	ctx *fiber.Ctx,
	destino string,
) error {
	tracingHeaders := propagation.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(ctx.UserContext(), tracingHeaders)
	for key, value := range tracingHeaders {
		if len(value) > 0 {
			ctx.Request().Header.Set(key, value[0])
		}
	}

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
