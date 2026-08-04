package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"practica2/backend/internal/config"
	"practica2/backend/internal/infrastructure/database"
	controlador "practica2/backend/internal/modules/auth/controller"
	rutas "practica2/backend/internal/modules/auth/routes"
	servicio "practica2/backend/internal/modules/auth/service"
	controladorProtegidas "practica2/backend/internal/modules/protected/controller"
	rutasProtegidas "practica2/backend/internal/modules/protected/routes"
	controladorUsuarios "practica2/backend/internal/modules/users/controller"
	repositorio "practica2/backend/internal/modules/users/repository"
	rutasUsuarios "practica2/backend/internal/modules/users/routes"
	servicioUsuarios "practica2/backend/internal/modules/users/service"
	"practica2/backend/internal/shared/cookies"
	seguridad "practica2/backend/internal/shared/crypto"
	token "practica2/backend/internal/shared/jwt"
	"practica2/backend/internal/shared/response"
)

func main() {
	config, err := configuracion.Cargar()
	if err != nil {
		log.Fatalf("error de configuracion: %v", err)
	}

	ctx := context.Background()
	pool, err := basedatos.Conectar(ctx, config.BaseDatos.URL)
	if err != nil {
		log.Fatalf("error al iniciar la base de datos: %v", err)
	}
	defer pool.Close()

	cifrador, err := seguridad.NuevoCifradorAES(config.Seguridad.ClaveAES)
	if err != nil {
		log.Fatalf("error al iniciar el cifrado: %v", err)
	}
	generadorCorreo, err := seguridad.NuevoGeneradorCorreo(config.Seguridad.ClaveCorreo)
	if err != nil {
		log.Fatalf("error al iniciar el indice de correo: %v", err)
	}
	vigencia := time.Duration(config.JWT.MinutosExpiracion) * time.Minute
	gracia := time.Duration(config.JWT.MinutosGracia) * time.Minute
	gestorToken, err := token.NuevoGestor(config.JWT.Secreto, vigencia, gracia, config.JWT.Emisor)
	if err != nil {
		log.Fatalf("error al iniciar JWT: %v", err)
	}
	gestorCookies, err := cookies.NuevoGestorSesion(config.Cookie.Segura, config.Cookie.SameSite, vigencia+gracia)
	if err != nil {
		log.Fatalf("error al iniciar cookies: %v", err)
	}

	repoUsuarios, err := repositorio.NuevoRepositorioPostgres(pool)
	if err != nil {
		log.Fatalf("error al iniciar usuarios: %v", err)
	}
	registro, err := servicio.NuevoServicioRegistro(repoUsuarios, repoUsuarios, cifrador, generadorCorreo)
	if err != nil {
		log.Fatalf("error al iniciar registro: %v", err)
	}
	login, err := servicio.NuevoServicioLogin(repoUsuarios, cifrador, gestorToken, generadorCorreo)
	if err != nil {
		log.Fatalf("error al iniciar login: %v", err)
	}
	auth, err := controlador.NuevaAutenticacion(registro, login, gestorCookies, gestorToken)
	if err != nil {
		log.Fatalf("error al iniciar autenticacion: %v", err)
	}
	consultaUsuarios, err := servicioUsuarios.NuevoServicioConsulta(repoUsuarios, cifrador)
	if err != nil {
		log.Fatalf("error al iniciar consulta de usuarios: %v", err)
	}
	estadoUsuarios, err := servicioUsuarios.NuevoServicioEstado(repoUsuarios, cifrador)
	if err != nil {
		log.Fatalf("error al iniciar estados de usuarios: %v", err)
	}
	usuarios, err := controladorUsuarios.NuevosUsuarios(consultaUsuarios, estadoUsuarios)
	if err != nil {
		log.Fatalf("error al iniciar administracion de usuarios: %v", err)
	}
	protegidas := controladorProtegidas.NuevaProtegidas()

	app := fiber.New(fiber.Config{
		AppName:               "API de Autenticacion - Practica 2",
		DisableStartupMessage: config.Aplicacion.Entorno == "test",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			estado := fiber.StatusInternalServerError
			codigo := "INTERNAL_ERROR"
			mensaje := "Ocurrio un error inesperado"

			var errorFiber *fiber.Error
			if errors.As(err, &errorFiber) {
				estado = errorFiber.Code
				mensaje = errorFiber.Message
				codigo = "HTTP_ERROR"
			}
			return respuesta.Fallo(ctx, estado, mensaje, codigo)
		},
	})

	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Aplicacion.URLFrontend,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	app.Get("/health", func(ctx *fiber.Ctx) error {
		ctxConsulta, cancelar := context.WithTimeout(ctx.Context(), 2*time.Second)
		defer cancelar()
		if err := pool.Ping(ctxConsulta); err != nil {
			return respuesta.Fallo(ctx, fiber.StatusServiceUnavailable, "Base de datos no disponible", "DATABASE_UNAVAILABLE")
		}
		return respuesta.Exitosa(ctx, fiber.StatusOK, "Servicio disponible", fiber.Map{"database": "connected"})
	})

	api := app.Group("/api/v1")
	rutas.RegistrarAutenticacion(api, auth, gestorCookies, gestorToken)
	rutasUsuarios.RegistrarUsuarios(api, usuarios, gestorCookies, gestorToken)
	rutasProtegidas.RegistrarProtegidas(api, protegidas, gestorCookies, gestorToken)

	go func() {
		direccion := fmt.Sprintf(":%d", config.Aplicacion.Puerto)
		if err := app.Listen(direccion); err != nil {
			log.Printf("servidor HTTP detenido: %v", err)
		}
	}()

	senal := make(chan os.Signal, 1)
	signal.Notify(senal, os.Interrupt, syscall.SIGTERM)
	<-senal

	ctxCierre, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := app.ShutdownWithContext(ctxCierre); err != nil {
		log.Printf("error al cerrar el servidor: %v", err)
	}
}
