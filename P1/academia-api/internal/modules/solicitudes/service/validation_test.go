package service

import (
	"context"
	"testing"

	"github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/modules/solicitudes/dto"
	errorsapp "github.com/Jeremy142OO4/Practicas-SA-B-202300644/P1/academia-api/internal/shared/errors"
)

func TestCrearExigeTodosLosCampos(
	t *testing.T,
) {
	t.Parallel()

	servicio := NewService(nil)

	_, err := servicio.Crear(
		context.Background(),
		dto.CrearSolicitudRequest{},
	)

	assertValidationError(t, err)
}

func TestActualizarExigeTodosLosCampos(
	t *testing.T,
) {
	t.Parallel()

	servicio := NewService(nil)

	_, err := servicio.Actualizar(
		context.Background(),
		"f677799e-a0c2-4a28-b018-d4901ae84dc8",
		dto.ActualizarSolicitudRequest{},
	)

	assertValidationError(t, err)
}

func TestActualizarEstadoExigeEstado(
	t *testing.T,
) {
	t.Parallel()

	servicio := NewService(nil)

	_, err := servicio.ActualizarEstado(
		context.Background(),
		"f677799e-a0c2-4a28-b018-d4901ae84dc8",
		dto.ActualizarEstadoRequest{},
	)

	assertValidationError(t, err)
}

func assertValidationError(
	t *testing.T,
	err error,
) {
	t.Helper()

	var appErr *errorsapp.AppError

	if !errorsapp.Es(err, &appErr) {
		t.Fatalf(
			"se esperaba AppError, se obtuvo %v",
			err,
		)
	}

	if appErr.Tipo != errorsapp.TipoValidacion {
		t.Fatalf(
			"tipo = %s, se esperaba %s",
			appErr.Tipo,
			errorsapp.TipoValidacion,
		)
	}
}
