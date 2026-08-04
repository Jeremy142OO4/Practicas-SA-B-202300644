package dto

type SolicitudRegistro struct {
	Nombre string `json:"name"`
	Correo string `json:"email"`
	Clave  string `json:"password"`
}

type UsuarioRegistrado struct {
	ID     string `json:"id"`
	Nombre string `json:"name"`
	Correo string `json:"email"`
	Rol    string `json:"role"`
	Estado string `json:"status"`
}
