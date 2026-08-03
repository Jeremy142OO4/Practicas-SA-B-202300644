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
