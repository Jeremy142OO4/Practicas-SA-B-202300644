package rutas

import (
	"github.com/gofiber/fiber/v2"

	controlador "practica2/backend/internal/modules/auth/controller"
	"practica2/backend/internal/modules/auth/middleware"
	"practica2/backend/internal/shared/cookies"
	token "practica2/backend/internal/shared/jwt"
)

func RegistrarAutenticacion(api fiber.Router, auth *controlador.Autenticacion, gestorCookies *cookies.GestorSesion, gestorToken *token.Gestor) {
	ruta := api.Group("/auth")
	ruta.Post("/register", auth.Registrar)
	ruta.Post("/login", auth.IniciarSesion)
	ruta.Post("/refresh", auth.Renovar)
	ruta.Post("/logout", auth.CerrarSesion)
	protegida := ruta.Group("", middleware.Autenticar(gestorCookies, gestorToken))
	protegida.Get("/session", auth.VerSesion)
	protegida.Get("/admin", middleware.Autorizar("ADMIN"), auth.VerAdministracion)
}
