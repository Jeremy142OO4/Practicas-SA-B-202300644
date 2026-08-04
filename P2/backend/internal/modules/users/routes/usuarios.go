package rutas

import (
	"github.com/gofiber/fiber/v2"

	"practica2/backend/internal/modules/auth/middleware"
	controlador "practica2/backend/internal/modules/users/controller"
	"practica2/backend/internal/shared/cookies"
	token "practica2/backend/internal/shared/jwt"
)

func RegistrarUsuarios(api fiber.Router, usuarios *controlador.Usuarios, gestorCookies *cookies.GestorSesion, gestorToken *token.Gestor) {
	ruta := api.Group("/users", middleware.Autenticar(gestorCookies, gestorToken), middleware.Autorizar("ADMIN"))
	ruta.Get("/", usuarios.Listar)
	ruta.Get("/:id", usuarios.Buscar)
	ruta.Patch("/:id/status", usuarios.CambiarEstado)
}
