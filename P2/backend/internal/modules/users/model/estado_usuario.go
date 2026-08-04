package modelo

type EstadoUsuario string

const (
	EstadoActivo    EstadoUsuario = "ACTIVE"
	EstadoInactivo  EstadoUsuario = "INACTIVE"
	EstadoBloqueado EstadoUsuario = "BLOCKED"
)

func (estado EstadoUsuario) EsValido() bool {
	switch estado {
	case EstadoActivo, EstadoInactivo, EstadoBloqueado:
		return true
	default:
		return false
	}
}
