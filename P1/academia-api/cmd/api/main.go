package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"

	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/config"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/infrastructure/database"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/controller"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/repository"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/routes"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/service"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/shared/response"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/shared/validator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"error de configuración: %v",
			err,
		)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	db, err := database.NewPostgresPool(
		ctx,
		cfg.Database,
	)
	if err != nil {
		log.Fatalf(
			"error de base de datos: %v",
			err,
		)
	}
	defer db.Close()

	solicitudRepository :=
		repository.NewPostgresRepository(db)

	solicitudService :=
		service.NewService(
			solicitudRepository,
		)

	solicitudController :=
		controller.NewController(
			solicitudService,
		)

	app := fiber.New(
		fiber.Config{
			AppName:     cfg.App.Name,
			JSONDecoder: validator.DecodificarJSONEstricto,
		},
	)

	api := app.Group("/api/v1")

	routes.Registrar(
		api,
		solicitudController,
	)

	app.Use(
		func(c fiber.Ctx) error {
			return response.Error(
				c,
				fiber.StatusNotFound,
				"ROUTE_NOT_FOUND",
				"ruta no encontrada",
			)
		},
	)

	address := fmt.Sprintf(
		"%s:%s",
		cfg.App.Host,
		cfg.App.Port,
	)

	go func() {
		log.Printf(
			"%s escuchando en http://%s",
			cfg.App.Name,
			address,
		)

		if err := app.Listen(address); err != nil {
			log.Printf(
				"servidor detenido: %v",
				err,
			)
		}
	}()

	<-ctx.Done()

	if err := app.Shutdown(); err != nil {
		log.Printf(
			"error al cerrar el servidor: %v",
			err,
		)
	}
}
