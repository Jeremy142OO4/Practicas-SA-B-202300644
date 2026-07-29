package errors

import "errors"

type Tipo string

const (
	TipoValidacion Tipo = "VALIDATION_ERROR"

	TipoNoEncontrado Tipo = "NOT_FOUND"

	TipoConflicto Tipo = "CONFLICT"

	TipoInterno Tipo = "INTERNAL_ERROR"
)

type AppError struct {
	Tipo    Tipo
	Mensaje string
	Causa   error
}

func (e *AppError) Error() string {
	return e.Mensaje
}

func (e *AppError) Unwrap() error {
	return e.Causa
}

func Nuevo(
	tipo Tipo,
	mensaje string,
	causa error,
) *AppError {
	return &AppError{
		Tipo:    tipo,
		Mensaje: mensaje,
		Causa:   causa,
	}
}

func Es(
	err error,
	target **AppError,
) bool {
	return errors.As(err, target)
}
