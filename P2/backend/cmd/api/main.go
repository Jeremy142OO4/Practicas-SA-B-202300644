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
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	ctx := context.Background()
	databasePool, err := database.Connect(ctx, appConfig.Database.URL)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	defer databasePool.Close()

	app := fiber.New(fiber.Config{
		AppName:               "Practice 2 Authentication API",
		DisableStartupMessage: appConfig.App.Environment == "test",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			code := "INTERNAL_ERROR"
			message := "An unexpected error occurred"

			var fiberError *fiber.Error
			if errors.As(err, &fiberError) {
				status = fiberError.Code
				message = fiberError.Message
				code = "HTTP_ERROR"
			}

			return response.Failure(ctx, status, message, code)
		},
	})

	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     appConfig.App.FrontendURL,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	app.Get("/health", func(ctx *fiber.Ctx) error {
		checkContext, cancel := context.WithTimeout(ctx.Context(), 2*time.Second)
		defer cancel()
		if err := databasePool.Ping(checkContext); err != nil {
			return response.Failure(ctx, fiber.StatusServiceUnavailable, "Database is unavailable", "DATABASE_UNAVAILABLE")
		}
		return response.Success(ctx, fiber.StatusOK, "Service is healthy", fiber.Map{"database": "connected"})
	})

	go func() {
		address := fmt.Sprintf(":%d", appConfig.App.Port)
		if err := app.Listen(address); err != nil {
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownContext); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
