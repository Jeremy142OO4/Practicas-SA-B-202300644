package models

import "time"

type Usuario struct {
	IDUsuario         string    `json:"id_usuario"`
	NameEncrypted     string    `json:"-"`
	EmailEncrypted    string    `json:"-"`
	EmailHash         string    `json:"-"`
	PasswordEncrypted string    `json:"-"`
	RoleID            int       `json:"role_id"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Rol Rol `json:"rol,omitempty"`
}

type RegistroRequest struct {
	Nombre     string `json:"nombre"`
	Correo     string `json:"correo"`
	Contrasena string `json:"password"`
}

type LoginRequest struct {
	Correo     string `json:"correo"`
	Contrasena string `json:"password"`
}

type UsuarioResponse struct {
	IDUsuario string `json:"id_usuario"`
	Nombre    string `json:"nombre"`
	Correo    string `json:"correo"`
	Rol       string `json:"rol"`
	Status    string `json:"status"`
}

type LoginResponse struct {
	Token   string          `json:"token"`
	Usuario UsuarioResponse `json:"usuario"`
}
