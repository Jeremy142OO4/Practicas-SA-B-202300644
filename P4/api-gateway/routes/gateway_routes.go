package routes

import (
	"github.com/gofiber/fiber/v2"

	"p4/api-gateway/config"
	"p4/api-gateway/handlers"
)

func RegistrarRutas(
	app *fiber.App,
	configuracion config.Configuracion,
	proxyHandler *handlers.ProxyHandler,
) {

	// ==========================================
	// AUTH-MS
	// ==========================================

	app.All(
		"/auth/*",
		func(ctx *fiber.Ctx) error {

			ruta := ctx.OriginalURL()

			destino :=
				configuracion.AuthURL +
					ruta

			return proxyHandler.Reenviar(
				ctx,
				destino,
			)
		},
	)

	// ==========================================
	// CATALOG-MS
	// ==========================================

	app.All(
		"/catalog/graphql",
		func(ctx *fiber.Ctx) error {

			destino :=
				configuracion.CatalogURL +
					"/graphql"

			return proxyHandler.Reenviar(
				ctx,
				destino,
			)
		},
	)

	// ==========================================
	// ORDERS-MS
	// ==========================================

	app.All(
		"/orders/graphql",
		func(ctx *fiber.Ctx) error {

			destino :=
				configuracion.OrdersURL +
					"/graphql"

			return proxyHandler.Reenviar(
				ctx,
				destino,
			)
		},
	)

	// ==========================================
	// PAYMENTS-MS
	// ==========================================

	app.All(
		"/payments/*",
		func(ctx *fiber.Ctx) error {

			ruta := ctx.Params("*")

			destino :=
				configuracion.PaymentsURL +
					"/pagos/" +
					ruta

			return proxyHandler.Reenviar(
				ctx,
				destino,
			)
		},
	)
}
