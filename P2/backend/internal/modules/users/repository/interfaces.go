package repositorio

import (
	"context"

	"practica2/backend/internal/modules/users/model"
)

type LectorUsuario interface {
	BuscarPorIndice(ctx context.Context, indice string) (*modelo.Usuario, error)
}

type ConsultorUsuarios interface {
	Listar(ctx context.Context) ([]modelo.Usuario, error)
	BuscarPorID(ctx context.Context, id string) (*modelo.Usuario, error)
}

type ActualizadorEstado interface {
	CambiarEstado(ctx context.Context, id string, estado modelo.EstadoUsuario) (*modelo.Usuario, error)
}

type EscritorUsuario interface {
	Crear(ctx context.Context, datos modelo.DatosUsuario) (*modelo.Usuario, error)
}

type RepositorioUsuario interface {
	LectorUsuario
	EscritorUsuario
	ConsultorUsuarios
	ActualizadorEstado
}
