package repositorio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"practica2/backend/internal/modules/users/model"
	seguridad "practica2/backend/internal/shared/crypto"
)

const indiceAdmin = "abf0788e2783e4a0caf57c4725c4488dbba89ff1891cc5823d8f41247ec667e6"

func TestBuscarUsuarioExistente(t *testing.T) {
	repositorio, pool := prepararRepositorio(t)
	usuario, err := repositorio.BuscarPorIndice(context.Background(), indiceAdmin)
	if err != nil {
		t.Fatalf("BuscarPorIndice() error = %v", err)
	}
	if usuario.Rol.Nombre != "ADMIN" || usuario.Estado != modelo.EstadoActivo {
		t.Fatalf("usuario inesperado: rol=%s estado=%s", usuario.Rol.Nombre, usuario.Estado)
	}
	if usuario.ID == "" || usuario.CorreoCifrado == "" {
		t.Fatal("usuario incompleto")
	}
	_ = pool
}

func TestBuscarUsuarioInexistente(t *testing.T) {
	repositorio, _ := prepararRepositorio(t)
	_, err := repositorio.BuscarPorIndice(context.Background(), fmt.Sprintf("%064d", 0))
	if !errors.Is(err, ErrUsuarioNoExiste) {
		t.Fatalf("BuscarPorIndice() error = %v", err)
	}
}

func TestCrearUsuario(t *testing.T) {
	repositorio, pool := prepararRepositorio(t)
	datos := crearDatos(t)
	defer borrarUsuario(t, pool, datos.IndiceCorreo)

	usuario, err := repositorio.Crear(context.Background(), datos)
	if err != nil {
		t.Fatalf("Crear() error = %v", err)
	}
	if usuario.Rol.Nombre != "CLIENTE" || usuario.Estado != modelo.EstadoActivo {
		t.Fatalf("usuario inesperado: rol=%s estado=%s", usuario.Rol.Nombre, usuario.Estado)
	}

	encontrado, err := repositorio.BuscarPorIndice(context.Background(), datos.IndiceCorreo)
	if err != nil || encontrado.ID != usuario.ID {
		t.Fatalf("BuscarPorIndice() = %v, error = %v", encontrado, err)
	}
}

func TestCrearCorreoDuplicado(t *testing.T) {
	repositorio, _ := prepararRepositorio(t)
	datos := crearDatos(t)
	datos.IndiceCorreo = indiceAdmin

	_, err := repositorio.Crear(context.Background(), datos)
	if !errors.Is(err, ErrCorreoDuplicado) {
		t.Fatalf("Crear() error = %v", err)
	}
}

func TestCrearConRolInexistente(t *testing.T) {
	repositorio, _ := prepararRepositorio(t)
	datos := crearDatos(t)
	datos.NombreRol = "DESCONOCIDO"

	_, err := repositorio.Crear(context.Background(), datos)
	if !errors.Is(err, ErrRolNoExiste) {
		t.Fatalf("Crear() error = %v", err)
	}
}

func prepararRepositorio(t *testing.T) (*RepositorioPostgres, *pgxpool.Pool) {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL no esta definida")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Fatalf("pool.Ping() error = %v", err)
	}
	t.Cleanup(pool.Close)
	repositorio, err := NuevoRepositorioPostgres(pool)
	if err != nil {
		t.Fatalf("NuevoRepositorioPostgres() error = %v", err)
	}
	return repositorio, pool
}

func crearDatos(t *testing.T) modelo.DatosUsuario {
	t.Helper()
	marca := time.Now().UnixNano()
	clave := make([]byte, seguridad.LongitudClaveAES)
	for i := range clave {
		clave[i] = 0x42
	}
	claveIndice := make([]byte, seguridad.LongitudClaveCorreo)
	for i := range claveIndice {
		claveIndice[i] = 0x24
	}
	cifrador, _ := seguridad.NuevoCifradorAES(clave)
	generador, _ := seguridad.NuevoGeneradorCorreo(claveIndice)
	correo := fmt.Sprintf("prueba-%d@local.test", marca)
	nombreCifrado, _ := cifrador.Cifrar("Usuario Temporal")
	correoCifrado, _ := cifrador.Cifrar(correo)
	claveCifrada, _ := cifrador.Cifrar("Temporal123!")
	return modelo.DatosUsuario{
		NombreCifrado: nombreCifrado,
		CorreoCifrado: correoCifrado,
		IndiceCorreo:  generador.Generar(correo),
		ClaveCifrada:  claveCifrada,
		NombreRol:     "CLIENTE",
		Estado:        modelo.EstadoActivo,
	}
}

func borrarUsuario(t *testing.T, pool *pgxpool.Pool, indice string) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), "DELETE FROM users WHERE email_hash = $1", indice); err != nil {
		t.Errorf("limpiar usuario temporal: %v", err)
	}
}
