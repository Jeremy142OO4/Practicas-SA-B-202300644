package cookies

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

const NombreSesion = "auth_token"

type GestorSesion struct {
	segura   bool
	sameSite string
	duracion time.Duration
}

func NuevoGestorSesion(segura bool, sameSite string, duracion time.Duration) (*GestorSesion, error) {
	if duracion <= 0 {
		return nil, errors.New("la duracion de la cookie debe ser positiva")
	}
	return &GestorSesion{segura: segura, sameSite: sameSite, duracion: duracion}, nil
}

func (gestor *GestorSesion) Guardar(ctx *fiber.Ctx, valor string) {
	ctx.Cookie(&fiber.Cookie{
		Name: NombreSesion, Value: valor, Path: "/",
		MaxAge:   int(gestor.duracion.Seconds()),
		HTTPOnly: true, Secure: gestor.segura, SameSite: gestor.sameSite,
	})
}

func (gestor *GestorSesion) Leer(ctx *fiber.Ctx) string {
	return ctx.Cookies(NombreSesion)
}

func (gestor *GestorSesion) Eliminar(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name: NombreSesion, Value: "", Path: "/", MaxAge: -1,
		Expires: time.Unix(0, 0), HTTPOnly: true,
		Secure: gestor.segura, SameSite: gestor.sameSite,
	})
}
