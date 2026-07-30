package dto

type ActualizarSolicitudRequest struct {
	Titulo *string `json:"titulo"`

	AreaSolicitante *string `json:"area_solicitante"`

	Prioridad *int16 `json:"prioridad"`

	CostoEstimado *float64 `json:"costo_estimado"`

	Estado *string `json:"estado"`
}
