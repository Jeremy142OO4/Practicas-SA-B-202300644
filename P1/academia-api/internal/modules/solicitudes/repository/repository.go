package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/model"
)

type Repository interface {
	ObtenerTodas(
		ctx context.Context,
	) ([]model.SolicitudOperativa, error)

	Crear(
		ctx context.Context,
		solicitud *model.SolicitudOperativa,
	) error

	Actualizar(
		ctx context.Context,
		solicitud *model.SolicitudOperativa,
	) error

	Eliminar(
		ctx context.Context,
		id uuid.UUID,
	) error

	ActualizarEstado(
		ctx context.Context,
		id uuid.UUID,
		estado model.Estado,
	) (*model.SolicitudOperativa, error)
}
