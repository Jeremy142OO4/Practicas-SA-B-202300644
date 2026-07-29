package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/tu-usuario/academia-api/internal/modules/solicitudes/model"
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
