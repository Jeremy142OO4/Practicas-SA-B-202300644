package servicio

import (
	"context"
	"errors"
	"testing"

	"practica2/backend/internal/modules/users/model"
)

const (
	idAdmin   = "00000000-0000-4000-8000-000000000001"
	idCliente = "00000000-0000-4000-8000-000000000002"
)

type repositorioPrueba struct {
	usuarios []modelo.Usuario
	ultimo   modelo.EstadoUsuario
}

func (repo *repositorioPrueba) Listar(context.Context) ([]modelo.Usuario, error) {
	return repo.usuarios, nil
}

func (repo *repositorioPrueba) BuscarPorID(_ context.Context, id string) (*modelo.Usuario, error) {
	for _, usuario := range repo.usuarios {
		if usuario.ID == id {
			return &usuario, nil
		}
	}
	return nil, errors.New("no existe")
}

func (repo *repositorioPrueba) CambiarEstado(_ context.Context, id string, estado modelo.EstadoUsuario) (*modelo.Usuario, error) {
	repo.ultimo = estado
	usuario := repo.usuarios[0]
	usuario.ID = id
	usuario.Estado = estado
	return &usuario, nil
}

type descifradorPrueba struct{}

func (descifradorPrueba) Descifrar(valor string) (string, error) { return "claro-" + valor, nil }

func TestListarUsuariosDescifraDatos(t *testing.T) {
	repo := &repositorioPrueba{usuarios: []modelo.Usuario{{
		ID: idCliente, NombreCifrado: "nombre", CorreoCifrado: "correo",
		Rol: modelo.Rol{Nombre: "CLIENTE"}, Estado: modelo.EstadoActivo,
	}}}
	servicio, _ := NuevoServicioConsulta(repo, descifradorPrueba{})
	usuarios, err := servicio.Listar(context.Background())
	if err != nil || len(usuarios) != 1 {
		t.Fatalf("usuarios = %+v, error = %v", usuarios, err)
	}
	if usuarios[0].Nombre != "claro-nombre" || usuarios[0].Correo != "claro-correo" {
		t.Fatalf("usuario no descifrado = %+v", usuarios[0])
	}
}

func TestCambiarEstadoValido(t *testing.T) {
	repo := &repositorioPrueba{usuarios: []modelo.Usuario{{NombreCifrado: "nombre", CorreoCifrado: "correo", Rol: modelo.Rol{Nombre: "CLIENTE"}}}}
	servicio, _ := NuevoServicioEstado(repo, descifradorPrueba{})
	usuario, err := servicio.Cambiar(context.Background(), idAdmin, idCliente, "blocked")
	if err != nil || usuario.Estado != "BLOCKED" || repo.ultimo != modelo.EstadoBloqueado {
		t.Fatalf("usuario = %+v, estado = %s, error = %v", usuario, repo.ultimo, err)
	}
}

func TestRechazaEstadoInvalido(t *testing.T) {
	servicio, _ := NuevoServicioEstado(&repositorioPrueba{}, descifradorPrueba{})
	if _, err := servicio.Cambiar(context.Background(), idAdmin, idCliente, "BORRADO"); !errors.Is(err, ErrEstadoInvalido) {
		t.Fatalf("error = %v", err)
	}
}

func TestRechazaCambioPropio(t *testing.T) {
	servicio, _ := NuevoServicioEstado(&repositorioPrueba{}, descifradorPrueba{})
	if _, err := servicio.Cambiar(context.Background(), idAdmin, idAdmin, "INACTIVE"); !errors.Is(err, ErrCambioPropio) {
		t.Fatalf("error = %v", err)
	}
}
