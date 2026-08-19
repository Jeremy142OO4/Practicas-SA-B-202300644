package models

type VideojuegoPlataforma struct {
	IDVideojuegoPlataforma string     `json:"id_videojuego_plataforma"`
	IDVideojuego           string     `json:"id_videojuego"`
	IDPlataforma           int        `json:"id_plataforma"`
	Precio                 float64    `json:"precio"`
	Stock                  int        `json:"stock"`
	Activo                 bool       `json:"activo"`
	Plataforma             Plataforma `json:"plataforma,omitempty"`
}

type CrearVideojuegoPlataformaInput struct {
	IDVideojuego string
	IDPlataforma int
	Precio       float64
	Stock        int
}
