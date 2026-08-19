package models

type Plataforma struct {
	IDPlataforma int    `json:"id_plataforma"`
	Nombre       string `json:"nombre"`
	Descripcion  string `json:"descripcion"`
}

type CrearPlataformaInput struct {
	Nombre      string
	Descripcion string
}
