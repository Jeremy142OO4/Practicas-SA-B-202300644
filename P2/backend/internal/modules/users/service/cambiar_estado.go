package servicio

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"practica2/backend/internal/modules/users/dto"
	"practica2/backend/internal/modules/users/model"
	repositorio "practica2/backend/internal/modules/users/repository"
	seguridad "practica2/backend/internal/shared/crypto"
)

type CambioEstados interface {
	Cambiar(ctx context.Context, administradorID string, usuarioID string, estado string) (*dto.Usuario, error)
}

type ServicioEstado struct {
	repositorio repositorio.ActualizadorEstado
	descifrador seguridad.Descifrador
}

func NuevoServicioEstado(repo repositorio.ActualizadorEstado, descifrador seguridad.Descifrador) (*ServicioEstado, error) {
	if repo == nil || descifrador == nil {
		return nil, errors.New("las dependencias de estado son obligatorias")
	}
	return &ServicioEstado{repositorio: repo, descifrador: descifrador}, nil
}

func (servicio *ServicioEstado) Cambiar(ctx context.Context, administradorID string, usuarioID string, valor string) (*dto.Usuario, error) {
	if _, err := uuid.Parse(usuarioID); err != nil {
		return nil, ErrIDInvalido
	}
	if administradorID == usuarioID {
		return nil, ErrCambioPropio
	}
	estado := modelo.EstadoUsuario(strings.ToUpper(strings.TrimSpace(valor)))
	if !estado.EsValido() {
		return nil, ErrEstadoInvalido
	}
	guardado, err := servicio.repositorio.CambiarEstado(ctx, usuarioID, estado)
	if err != nil {
		return nil, fmt.Errorf("actualizar usuario: %w", err)
	}
	usuario, err := convertirUsuario(*guardado, servicio.descifrador)
	return &usuario, err
}

var _ CambioEstados = (*ServicioEstado)(nil)
