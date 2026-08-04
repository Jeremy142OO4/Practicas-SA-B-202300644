package servicio

import "errors"

var (
	ErrIDInvalido     = errors.New("el identificador no es valido")
	ErrEstadoInvalido = errors.New("el estado debe ser ACTIVE, INACTIVE o BLOCKED")
	ErrCambioPropio   = errors.New("no puedes cambiar el estado de tu propio usuario")
)
