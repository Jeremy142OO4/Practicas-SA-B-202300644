package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConectarBaseDatos() (*pgxpool.Pool, error) {
	conexion := os.Getenv("DATABASE_URL")

	if conexion == "" {
		modoSSL := os.Getenv("DB_SSLMODE")
		if modoSSL == "" {
			modoSSL = "disable"
		}
		conexion = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"), modoSSL,
		)
	}

	bd, err := pgxpool.New(
		context.Background(),
		conexion,
	)

	if err != nil {
		return nil, err
	}

	if err := bd.Ping(context.Background()); err != nil {
		bd.Close()
		return nil, err
	}

	return bd, nil
}
