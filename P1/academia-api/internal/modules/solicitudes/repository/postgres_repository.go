package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tu-usuario/academia-api/internal/modules/solicitudes/model"
	errorsapp "github.com/tu-usuario/academia-api/internal/shared/errors"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(
	db *pgxpool.Pool,
) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) ObtenerTodas(
	ctx context.Context,
) ([]model.SolicitudOperativa, error) {
	const query = `
		SELECT
			id,
			titulo,
			area_solicitante,
			prioridad,
			costo_estimado::double precision,
			estado,
			created_at,
			updated_at
		FROM public.solicitudes_operativas
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"consultar solicitudes: %w",
			err,
		)
	}
	defer rows.Close()

	solicitudes := make(
		[]model.SolicitudOperativa,
		0,
	)

	for rows.Next() {
		var solicitud model.SolicitudOperativa

		err := rows.Scan(
			&solicitud.ID,
			&solicitud.Titulo,
			&solicitud.AreaSolicitante,
			&solicitud.Prioridad,
			&solicitud.CostoEstimado,
			&solicitud.Estado,
			&solicitud.CreatedAt,
			&solicitud.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"leer solicitud: %w",
				err,
			)
		}

		solicitudes = append(
			solicitudes,
			solicitud,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"recorrer solicitudes: %w",
			err,
		)
	}

	return solicitudes, nil
}

func (r *postgresRepository) Crear(
	ctx context.Context,
	solicitud *model.SolicitudOperativa,
) error {
	const query = `
		INSERT INTO public.solicitudes_operativas (
			titulo,
			area_solicitante,
			prioridad,
			costo_estimado,
			estado
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		solicitud.Titulo,
		solicitud.AreaSolicitante,
		solicitud.Prioridad,
		solicitud.CostoEstimado,
		solicitud.Estado,
	).Scan(
		&solicitud.ID,
		&solicitud.CreatedAt,
		&solicitud.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"crear solicitud: %w",
			err,
		)
	}

	return nil
}

func (r *postgresRepository) Actualizar(
	ctx context.Context,
	solicitud *model.SolicitudOperativa,
) error {
	const query = `
		UPDATE public.solicitudes_operativas
		SET
			titulo = $2,
			area_solicitante = $3,
			prioridad = $4,
			costo_estimado = $5,
			estado = $6
		WHERE id = $1
		RETURNING
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		solicitud.ID,
		solicitud.Titulo,
		solicitud.AreaSolicitante,
		solicitud.Prioridad,
		solicitud.CostoEstimado,
		solicitud.Estado,
	).Scan(
		&solicitud.CreatedAt,
		&solicitud.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return errorsapp.Nuevo(
			errorsapp.TipoNoEncontrado,
			"solicitud no encontrada",
			err,
		)
	}

	if err != nil {
		return fmt.Errorf(
			"actualizar solicitud: %w",
			err,
		)
	}

	return nil
}

func (r *postgresRepository) Eliminar(
	ctx context.Context,
	id uuid.UUID,
) error {
	const query = `
		DELETE FROM public.solicitudes_operativas
		WHERE id = $1
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"eliminar solicitud: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return errorsapp.Nuevo(
			errorsapp.TipoNoEncontrado,
			"solicitud no encontrada",
			pgx.ErrNoRows,
		)
	}

	return nil
}

func (r *postgresRepository) ActualizarEstado(
	ctx context.Context,
	id uuid.UUID,
	estado model.Estado,
) (*model.SolicitudOperativa, error) {
	const query = `
		UPDATE public.solicitudes_operativas
		SET estado = $2
		WHERE id = $1
		RETURNING
			id,
			titulo,
			area_solicitante,
			prioridad,
			costo_estimado::double precision,
			estado,
			created_at,
			updated_at
	`

	var solicitud model.SolicitudOperativa

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		estado,
	).Scan(
		&solicitud.ID,
		&solicitud.Titulo,
		&solicitud.AreaSolicitante,
		&solicitud.Prioridad,
		&solicitud.CostoEstimado,
		&solicitud.Estado,
		&solicitud.CreatedAt,
		&solicitud.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorsapp.Nuevo(
			errorsapp.TipoNoEncontrado,
			"solicitud no encontrada",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"actualizar estado: %w",
			err,
		)
	}

	return &solicitud, nil
}
