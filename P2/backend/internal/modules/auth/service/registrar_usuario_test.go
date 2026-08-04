package servicio

import (
	"context"
	"errors"
	"testing"

	"practica2/backend/internal/modules/auth/dto"
	"practica2/backend/internal/modules/users/model"
	repositorio "practica2/backend/internal/modules/users/repository"
)

func TestRegistrarUsuario(t *testing.T) {
	repo := &repositorioPrueba{errorBuscar: repositorio.ErrUsuarioNoExiste}
	servicio := crearServicio(t, repo)

	usuario, err := servicio.Registrar(context.Background(), dto.SolicitudRegistro{
		Nombre: " Usuario Prueba ", Correo: " USUARIO@EJEMPLO.COM ", Clave: "Segura123!",
	})
	if err != nil {
		t.Fatalf("Registrar() error = %v", err)
	}
	if usuario.Correo != "usuario@ejemplo.com" || usuario.Rol != "CLIENTE" {
		t.Fatalf("usuario inesperado: %+v", usuario)
	}
	if repo.datos.NombreRol != "CLIENTE" || repo.datos.Estado != modelo.EstadoActivo {
		t.Fatalf("datos inesperados: %+v", repo.datos)
	}
}

func TestRegistrarCorreoDuplicado(t *testing.T) {
	repo := &repositorioPrueba{usuario: &modelo.Usuario{ID: "existente"}}
	servicio := crearServicio(t, repo)
	_, err := servicio.Registrar(context.Background(), dto.SolicitudRegistro{
		Nombre: "Usuario", Correo: "usuario@ejemplo.com", Clave: "Segura123!",
	})
	if !errors.Is(err, repositorio.ErrCorreoDuplicado) {
		t.Fatalf("Registrar() error = %v", err)
	}
}

func TestRegistrarDatosInvalidos(t *testing.T) {
	servicio := crearServicio(t, &repositorioPrueba{})
	_, err := servicio.Registrar(context.Background(), dto.SolicitudRegistro{
		Nombre: "A", Correo: "incorrecto", Clave: "debil",
	})
	if err == nil {
		t.Fatal("Registrar() debio rechazar los datos")
	}
}

type repositorioPrueba struct {
	usuario     *modelo.Usuario
	errorBuscar error
	datos       modelo.DatosUsuario
}

func (repo *repositorioPrueba) BuscarPorIndice(context.Context, string) (*modelo.Usuario, error) {
	return repo.usuario, repo.errorBuscar
}

func (repo *repositorioPrueba) Crear(_ context.Context, datos modelo.DatosUsuario) (*modelo.Usuario, error) {
	repo.datos = datos
	return &modelo.Usuario{
		ID: "nuevo", Rol: modelo.Rol{Nombre: "CLIENTE"}, Estado: modelo.EstadoActivo,
	}, nil
}

type cifradorPrueba struct{}

func (cifradorPrueba) Cifrar(valor string) (string, error) { return "cifrado:" + valor, nil }

type indicePrueba struct{}

func (indicePrueba) Generar(correo string) string { return "indice:" + correo }

func crearServicio(t *testing.T, repo *repositorioPrueba) *ServicioRegistro {
	t.Helper()
	servicio, err := NuevoServicioRegistro(repo, repo, cifradorPrueba{}, indicePrueba{})
	if err != nil {
		t.Fatalf("NuevoServicioRegistro() error = %v", err)
	}
	return servicio
}
