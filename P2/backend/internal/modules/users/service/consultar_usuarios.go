package servicio

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"practica2/backend/internal/modules/users/dto"
	repositorio "practica2/backend/internal/modules/users/repository"
	seguridad "practica2/backend/internal/shared/crypto"
)

type ConsultaUsuarios interface {
	Listar(ctx context.Context) ([]dto.Usuario, error)
	Buscar(ctx context.Context, id string) (*dto.Usuario, error)
}

type ServicioConsulta struct {
	repositorio repositorio.ConsultorUsuarios
	descifrador seguridad.Descifrador
}

func NuevoServicioConsulta(repo repositorio.ConsultorUsuarios, descifrador seguridad.Descifrador) (*ServicioConsulta, error) {
	if repo == nil || descifrador == nil {
		return nil, errors.New("las dependencias de consulta son obligatorias")
	}
	return &ServicioConsulta{repositorio: repo, descifrador: descifrador}, nil
}

func (servicio *ServicioConsulta) Listar(ctx context.Context) ([]dto.Usuario, error) {
	guardados, err := servicio.repositorio.Listar(ctx)
	if err != nil {
		return nil, err
	}
	usuarios := make([]dto.Usuario, 0, len(guardados))
	for _, guardado := range guardados {
		usuario, err := convertirUsuario(guardado, servicio.descifrador)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (servicio *ServicioConsulta) Buscar(ctx context.Context, id string) (*dto.Usuario, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrIDInvalido
	}
	guardado, err := servicio.repositorio.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("consultar usuario: %w", err)
	}
	usuario, err := convertirUsuario(*guardado, servicio.descifrador)
	return &usuario, err
}

var _ ConsultaUsuarios = (*ServicioConsulta)(nil)
