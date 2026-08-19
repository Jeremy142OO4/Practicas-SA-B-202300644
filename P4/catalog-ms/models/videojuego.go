package models

import "time"

type Videojuego struct {
	IDVideojuego     string    `json:"id_videojuego"`
	Titulo           string    `json:"titulo"`
	Descripcion      string    `json:"descripcion"`
	IDCategoria      int       `json:"id_categoria"`
	Desarrollador    string    `json:"desarrollador"`
	FechaLanzamiento time.Time `json:"fecha_lanzamiento"`
	Estado           string    `json:"estado"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CrearVideojuegoInput struct {
	Titulo           string
	Descripcion      string
	IDCategoria      int
	Desarrollador    string
	FechaLanzamiento string
}
