package servicio

import (
	"context"
	"errors"
	"testing"

	"practica2/backend/internal/modules/auth/dto"
	"practica2/backend/internal/modules/users/model"
	repositorio "practica2/backend/internal/modules/users/repository"
)

func TestIniciarSesion(t *testing.T) {
	lector := &lectorLogin{usuario: usuarioLogin()}
	servicio := crearLogin(t, lector)
	resultado, err := servicio.Iniciar(context.Background(), dto.SolicitudLogin{
		Correo: " ADMIN@PRACTICA.LOCAL ", Clave: "Admin123!",
	})
	if err != nil {
		t.Fatalf("Iniciar() error = %v", err)
	}
	if resultado.Token != "token-firmado" || resultado.Usuario.Rol != "ADMIN" {
		t.Fatalf("resultado inesperado: %+v", resultado)
	}
}

func TestLoginClaveIncorrecta(t *testing.T) {
	servicio := crearLogin(t, &lectorLogin{usuario: usuarioLogin()})
	_, err := servicio.Iniciar(context.Background(), dto.SolicitudLogin{
		Correo: "admin@practica.local", Clave: "Incorrecta123!",
	})
	if !errors.Is(err, ErrCredenciales) {
		t.Fatalf("Iniciar() error = %v", err)
	}
}

func TestLoginUsuarioInexistente(t *testing.T) {
	servicio := crearLogin(t, &lectorLogin{err: repositorio.ErrUsuarioNoExiste})
	_, err := servicio.Iniciar(context.Background(), dto.SolicitudLogin{
		Correo: "nadie@local.test", Clave: "Incorrecta123!",
	})
	if !errors.Is(err, ErrCredenciales) {
		t.Fatalf("Iniciar() error = %v", err)
	}
}

type lectorLogin struct {
	usuario *modelo.Usuario
	err     error
}

func (lector *lectorLogin) BuscarPorIndice(context.Context, string) (*modelo.Usuario, error) {
	return lector.usuario, lector.err
}

type descifradorLogin struct{}

func (descifradorLogin) Descifrar(valor string) (string, error) {
	return map[string]string{
		"nombre-cifrado": "Administrador Pruebas",
		"correo-cifrado": "admin@practica.local",
		"clave-cifrada":  "Admin123!",
	}[valor], nil
}

type tokenLogin struct{}

func (tokenLogin) Generar(string, string) (string, error) { return "token-firmado", nil }

type indiceLogin struct{}

func (indiceLogin) Generar(string) string { return "indice" }

func usuarioLogin() *modelo.Usuario {
	return &modelo.Usuario{
		ID: "usuario-1", NombreCifrado: "nombre-cifrado",
		CorreoCifrado: "correo-cifrado", ClaveCifrada: "clave-cifrada",
		Rol: modelo.Rol{Nombre: "ADMIN"}, Estado: modelo.EstadoActivo,
	}
}

func crearLogin(t *testing.T, lector repositorio.LectorUsuario) *ServicioLogin {
	t.Helper()
	servicio, err := NuevoServicioLogin(lector, descifradorLogin{}, tokenLogin{}, indiceLogin{})
	if err != nil {
		t.Fatalf("NuevoServicioLogin() error = %v", err)
	}
	return servicio
}
