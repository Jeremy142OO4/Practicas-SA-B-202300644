package rutas

import (
	"github.com/gofiber/fiber/v2"

	"practica2/backend/internal/modules/auth/middleware"
	controlador "practica2/backend/internal/modules/protected/controller"
	"practica2/backend/internal/shared/cookies"
	token "practica2/backend/internal/shared/jwt"
)

func RegistrarProtegidas(api fiber.Router, rutas *controlador.Protegidas, gestorCookies *cookies.GestorSesion, gestorToken *token.Gestor) {
	protegida := api.Group("/protected", middleware.Autenticar(gestorCookies, gestorToken))
	protegida.Get("/admin", middleware.Autorizar("ADMIN"), rutas.Administrador)
	protegida.Get("/shared", middleware.Autorizar("ADMIN", "CLIENTE"), rutas.Compartida)
}
