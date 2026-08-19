package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConectarBaseDatos() (*pgxpool.Pool, error) {

	host := os.Getenv("DB_HOST")
	puerto := os.Getenv("DB_PORT")
	usuario := os.Getenv("DB_USER")
	contrasena := os.Getenv("DB_PASSWORD")
	nombreBD := os.Getenv("DB_NAME")

	conexion := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		usuario,
		contrasena,
		host,
		puerto,
		nombreBD,
	)

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
