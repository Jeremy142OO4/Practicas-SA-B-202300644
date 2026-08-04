package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"practica2/backend/internal/shared/cookies"
	token "practica2/backend/internal/shared/jwt"
)

const secretoPrueba = "secreto-de-prueba-con-mas-de-32-caracteres"

func TestRutaRequiereSesion(t *testing.T) {
	app, _, _ := crearAplicacion(t)
	respuesta, err := app.Test(httptest.NewRequest("GET", "/protegida", nil))
	if err != nil || respuesta.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("estado = %d, error = %v", respuesta.StatusCode, err)
	}
}

func TestAdministradorPuedeEntrar(t *testing.T) {
	app, gestorToken, _ := crearAplicacion(t)
	valor, _ := gestorToken.Generar("usuario-1", "ADMIN")
	solicitud := httptest.NewRequest("GET", "/admin", nil)
	solicitud.AddCookie(&http.Cookie{Name: cookies.NombreSesion, Value: valor})
	respuesta, err := app.Test(solicitud)
	if err != nil || respuesta.StatusCode != fiber.StatusOK {
		t.Fatalf("estado = %d, error = %v", respuesta.StatusCode, err)
	}
}

func TestClienteNoPuedeEntrarComoAdministrador(t *testing.T) {
	app, gestorToken, _ := crearAplicacion(t)
	valor, _ := gestorToken.Generar("usuario-2", "CLIENTE")
	solicitud := httptest.NewRequest("GET", "/admin", nil)
	solicitud.AddCookie(&http.Cookie{Name: cookies.NombreSesion, Value: valor})
	respuesta, err := app.Test(solicitud)
	if err != nil || respuesta.StatusCode != fiber.StatusForbidden {
		t.Fatalf("estado = %d, error = %v", respuesta.StatusCode, err)
	}
}

func crearAplicacion(t *testing.T) (*fiber.App, *token.Gestor, *cookies.GestorSesion) {
	t.Helper()
	gestorToken, err := token.NuevoGestor(secretoPrueba, time.Minute, time.Minute, "practica2-api")
	if err != nil {
		t.Fatal(err)
	}
	gestorCookies, err := cookies.NuevoGestorSesion(false, "lax", 2*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	app := fiber.New()
	app.Get("/protegida", Autenticar(gestorCookies, gestorToken), func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
	app.Get("/admin", Autenticar(gestorCookies, gestorToken), Autorizar("ADMIN"), func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
	return app, gestorToken, gestorCookies
}
