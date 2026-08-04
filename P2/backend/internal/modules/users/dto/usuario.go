package dto

import "time"

type Usuario struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"name"`
	Correo         string    `json:"email"`
	Rol            string    `json:"role"`
	Estado         string    `json:"status"`
	FechaCreacion  time.Time `json:"createdAt"`
	FechaActualiza time.Time `json:"updatedAt"`
}

type CambioEstado struct {
	Estado string `json:"status"`
}
