package models

type Categoria struct {
	IDCategoria int    `json:"id_categoria"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type CrearCategoriaInput struct {
	Nombre      string
	Descripcion string
}
