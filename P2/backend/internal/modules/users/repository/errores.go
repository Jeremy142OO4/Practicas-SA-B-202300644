package repositorio

import "errors"

var (
	ErrUsuarioNoExiste = errors.New("usuario no encontrado")
	ErrCorreoDuplicado = errors.New("correo ya registrado")
	ErrRolNoExiste     = errors.New("rol no encontrado")
)
