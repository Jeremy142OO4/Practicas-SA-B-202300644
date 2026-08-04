package dto

type SolicitudLogin struct {
	Correo string `json:"email"`
	Clave  string `json:"password"`
}

type UsuarioSesion struct {
	ID     string `json:"id"`
	Nombre string `json:"name"`
	Correo string `json:"email"`
	Rol    string `json:"role"`
	Estado string `json:"status"`
}
