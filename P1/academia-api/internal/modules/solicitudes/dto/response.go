package dto

import "time"

type SolicitudResponse struct {
	ID string `json:"id"`

	Titulo string `json:"titulo"`

	AreaSolicitante string `json:"area_solicitante"`

	Prioridad int16 `json:"prioridad"`

	CostoEstimado float64 `json:"costo_estimado"`

	Estado string `json:"estado"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
