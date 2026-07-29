package model

import (
	"time"

	"github.com/google/uuid"
)

type Estado string

const (
	EstadoRegistrada Estado = "registrada"
	EstadoEnProceso  Estado = "en_proceso"
	EstadoFinalizada Estado = "finalizada"
)

type SolicitudOperativa struct {
	ID              uuid.UUID
	Titulo          string
	AreaSolicitante string
	Prioridad       int16
	CostoEstimado   float64
	Estado          Estado
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (e Estado) EsValido() bool {
	switch e {
	case EstadoRegistrada,
		EstadoEnProceso,
		EstadoFinalizada:
		return true
	default:
		return false
	}
}
