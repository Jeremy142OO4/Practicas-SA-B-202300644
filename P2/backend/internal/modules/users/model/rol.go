package modelo

type Rol struct {
	ID          int16  `json:"id"`
	Nombre      string `json:"name"`
	Descripcion string `json:"description"`
}
