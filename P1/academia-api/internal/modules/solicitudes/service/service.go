package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/dto"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/model"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/repository"
	errorsapp "github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/shared/errors"
	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/shared/validator"
)

type Service interface {
	ObtenerTodas(
		ctx context.Context,
	) ([]dto.SolicitudResponse, error)

	Crear(
		ctx context.Context,
		input dto.CrearSolicitudRequest,
	) (*dto.SolicitudResponse, error)

	Actualizar(
		ctx context.Context,
		id string,
		input dto.ActualizarSolicitudRequest,
	) (*dto.SolicitudResponse, error)

	Eliminar(
		ctx context.Context,
		id string,
	) error

	ActualizarEstado(
		ctx context.Context,
		id string,
		input dto.ActualizarEstadoRequest,
	) (*dto.SolicitudResponse, error)
}

type solicitudService struct {
	repository repository.Repository
}

func NewService(
	repository repository.Repository,
) Service {
	return &solicitudService{
		repository: repository,
	}
}

func (s *solicitudService) ObtenerTodas(
	ctx context.Context,
) ([]dto.SolicitudResponse, error) {
	solicitudes, err :=
		s.repository.ObtenerTodas(ctx)

	if err != nil {
		return nil, errorsapp.Nuevo(
			errorsapp.TipoInterno,
			"no fue posible obtener las solicitudes",
			err,
		)
	}

	responses := make(
		[]dto.SolicitudResponse,
		0,
		len(solicitudes),
	)

	for _, solicitud := range solicitudes {
		responses = append(
			responses,
			convertirAResponse(solicitud),
		)
	}

	return responses, nil
}

func (s *solicitudService) Crear(
	ctx context.Context,
	input dto.CrearSolicitudRequest,
) (*dto.SolicitudResponse, error) {
	err := validarDatos(
		input.Titulo,
		input.AreaSolicitante,
		input.Prioridad,
		input.CostoEstimado,
	)
	if err != nil {
		return nil, err
	}

	solicitud := model.SolicitudOperativa{
		Titulo: validator.LimpiarTexto(
			input.Titulo,
		),
		AreaSolicitante: validator.LimpiarTexto(
			input.AreaSolicitante,
		),
		Prioridad:     input.Prioridad,
		CostoEstimado: input.CostoEstimado,
		Estado:        model.EstadoRegistrada,
	}

	if err := s.repository.Crear(
		ctx,
		&solicitud,
	); err != nil {
		return nil, errorsapp.Nuevo(
			errorsapp.TipoInterno,
			"no fue posible crear la solicitud",
			err,
		)
	}

	resultado := convertirAResponse(solicitud)

	return &resultado, nil
}

func (s *solicitudService) Actualizar(
	ctx context.Context,
	id string,
	input dto.ActualizarSolicitudRequest,
) (*dto.SolicitudResponse, error) {
	uuidValue, err := validarUUID(id)
	if err != nil {
		return nil, err
	}

	err = validarDatos(
		input.Titulo,
		input.AreaSolicitante,
		input.Prioridad,
		input.CostoEstimado,
	)
	if err != nil {
		return nil, err
	}

	estado := model.Estado(input.Estado)

	if !estado.EsValido() {
		return nil, errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			"estado inválido",
			nil,
		)
	}

	solicitud := model.SolicitudOperativa{
		ID: uuidValue,
		Titulo: validator.LimpiarTexto(
			input.Titulo,
		),
		AreaSolicitante: validator.LimpiarTexto(
			input.AreaSolicitante,
		),
		Prioridad:     input.Prioridad,
		CostoEstimado: input.CostoEstimado,
		Estado:        estado,
	}

	if err := s.repository.Actualizar(
		ctx,
		&solicitud,
	); err != nil {
		return nil, normalizarError(
			err,
			"no fue posible actualizar la solicitud",
		)
	}

	resultado := convertirAResponse(solicitud)

	return &resultado, nil
}

func (s *solicitudService) Eliminar(
	ctx context.Context,
	id string,
) error {
	uuidValue, err := validarUUID(id)
	if err != nil {
		return err
	}

	if err := s.repository.Eliminar(
		ctx,
		uuidValue,
	); err != nil {
		return normalizarError(
			err,
			"no fue posible eliminar la solicitud",
		)
	}

	return nil
}

func (s *solicitudService) ActualizarEstado(
	ctx context.Context,
	id string,
	input dto.ActualizarEstadoRequest,
) (*dto.SolicitudResponse, error) {
	uuidValue, err := validarUUID(id)
	if err != nil {
		return nil, err
	}

	estado := model.Estado(input.Estado)

	if !estado.EsValido() {
		return nil, errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			"estado debe ser registrada, "+
				"en_proceso o finalizada",
			nil,
		)
	}

	solicitud, err :=
		s.repository.ActualizarEstado(
			ctx,
			uuidValue,
			estado,
		)

	if err != nil {
		return nil, normalizarError(
			err,
			"no fue posible actualizar el estado",
		)
	}

	resultado := convertirAResponse(*solicitud)

	return &resultado, nil
}

func validarDatos(
	titulo string,
	area string,
	prioridad int16,
	costo float64,
) error {
	if message := validator.ValidarTextoRequerido(
		titulo,
		"titulo",
		150,
	); message != "" {
		return errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			message,
			nil,
		)
	}

	if message := validator.ValidarTextoRequerido(
		area,
		"area_solicitante",
		100,
	); message != "" {
		return errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			message,
			nil,
		)
	}

	if prioridad < 1 || prioridad > 5 {
		return errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			"prioridad debe estar entre 1 y 5",
			nil,
		)
	}

	if costo < 0 {
		return errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			"costo_estimado no puede ser negativo",
			nil,
		)
	}

	return nil
}

func validarUUID(
	id string,
) (uuid.UUID, error) {
	value, err := uuid.Parse(id)

	if err != nil {
		return uuid.Nil, errorsapp.Nuevo(
			errorsapp.TipoValidacion,
			"id no es un UUID válido",
			err,
		)
	}

	return value, nil
}

func convertirAResponse(
	solicitud model.SolicitudOperativa,
) dto.SolicitudResponse {
	return dto.SolicitudResponse{
		ID:              solicitud.ID.String(),
		Titulo:          solicitud.Titulo,
		AreaSolicitante: solicitud.AreaSolicitante,
		Prioridad:       solicitud.Prioridad,
		CostoEstimado:   solicitud.CostoEstimado,
		Estado:          string(solicitud.Estado),
		CreatedAt:       solicitud.CreatedAt,
		UpdatedAt:       solicitud.UpdatedAt,
	}
}

func normalizarError(
	err error,
	mensajeInterno string,
) error {
	var appErr *errorsapp.AppError

	if errorsapp.Es(err, &appErr) {
		return appErr
	}

	return errorsapp.Nuevo(
		errorsapp.TipoInterno,
		mensajeInterno,
		fmt.Errorf(
			"operación de repositorio: %w",
			err,
		),
	)
}
