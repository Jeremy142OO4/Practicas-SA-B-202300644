package token

import (
	"errors"
	"testing"
	"time"
)

const secretoPrueba = "secreto-de-prueba-con-mas-de-32-caracteres"

func TestGenerarYValidar(t *testing.T) {
	gestor := crearGestor(t)
	valor, err := gestor.Generar("usuario-1", "ADMIN")
	if err != nil {
		t.Fatalf("Generar() error = %v", err)
	}
	datos, err := gestor.Validar(valor)
	if err != nil || datos.ID != "usuario-1" || datos.Rol != "ADMIN" {
		t.Fatalf("Validar() datos = %+v, error = %v", datos, err)
	}
}

func TestTokenExpiraEn60Segundos(t *testing.T) {
	gestor := crearGestor(t)
	inicio := time.Unix(1_700_000_000, 0)
	gestor.reloj = func() time.Time { return inicio }
	valor, _ := gestor.Generar("usuario-1", "CLIENTE")
	datos, _ := gestor.Validar(valor)
	if datos.ExpiraEn-inicio.Unix() != 60 {
		t.Fatalf("vigencia = %d segundos", datos.ExpiraEn-inicio.Unix())
	}

	gestor.reloj = func() time.Time { return inicio.Add(61 * time.Second) }
	if _, err := gestor.Validar(valor); !errors.Is(err, ErrTokenExpirado) {
		t.Fatalf("Validar() error = %v", err)
	}
}

func TestRenovacionDentroDeGracia(t *testing.T) {
	gestor := crearGestor(t)
	inicio := time.Unix(1_700_000_000, 0)
	gestor.reloj = func() time.Time { return inicio }
	valor, _ := gestor.Generar("usuario-1", "CLIENTE")

	gestor.reloj = func() time.Time { return inicio.Add(90 * time.Second) }
	datos, err := gestor.ValidarRenovacion(valor)
	if err != nil || datos.ID != "usuario-1" {
		t.Fatalf("ValidarRenovacion() datos = %+v, error = %v", datos, err)
	}
}

func TestRenovacionFueraDeGracia(t *testing.T) {
	gestor := crearGestor(t)
	inicio := time.Unix(1_700_000_000, 0)
	gestor.reloj = func() time.Time { return inicio }
	valor, _ := gestor.Generar("usuario-1", "CLIENTE")

	gestor.reloj = func() time.Time { return inicio.Add(121 * time.Second) }
	if _, err := gestor.ValidarRenovacion(valor); !errors.Is(err, ErrFueraDeGracia) {
		t.Fatalf("ValidarRenovacion() error = %v", err)
	}
}

func TestTokenManipulado(t *testing.T) {
	gestor := crearGestor(t)
	valor, _ := gestor.Generar("usuario-1", "CLIENTE")
	valor = valor[:len(valor)-1] + "x"
	if _, err := gestor.Validar(valor); !errors.Is(err, ErrTokenInvalido) {
		t.Fatalf("Validar() error = %v", err)
	}
}

func crearGestor(t *testing.T) *Gestor {
	t.Helper()
	gestor, err := NuevoGestor(secretoPrueba, time.Minute, time.Minute, "practica2-api")
	if err != nil {
		t.Fatalf("NuevoGestor() error = %v", err)
	}
	return gestor
}
