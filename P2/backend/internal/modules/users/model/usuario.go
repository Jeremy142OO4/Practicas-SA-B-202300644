package modelo

import "time"

type Usuario struct {
	ID             string        `json:"id"`
	NombreCifrado  string        `json:"-"`
	CorreoCifrado  string        `json:"-"`
	IndiceCorreo   string        `json:"-"`
	ClaveCifrada   string        `json:"-"`
	Rol            Rol           `json:"role"`
	Estado         EstadoUsuario `json:"status"`
	FechaCreacion  time.Time     `json:"createdAt"`
	FechaActualiza time.Time     `json:"updatedAt"`
}

type DatosUsuario struct {
	NombreCifrado string
	CorreoCifrado string
	IndiceCorreo  string
	ClaveCifrada  string
	NombreRol     string
	Estado        EstadoUsuario
}
