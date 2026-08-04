package validacion

import (
	"errors"
	"testing"
)

func TestValidarRegistro(t *testing.T) {
	casos := []struct {
		nombre string
		correo string
		clave  string
		espera error
	}{
		{"Usuario Valido", "usuario@ejemplo.com", "Segura123!", nil},
		{"A", "usuario@ejemplo.com", "Segura123!", ErrNombreInvalido},
		{"Usuario", "correo-invalido", "Segura123!", ErrCorreoInvalido},
		{"Usuario", "usuario@ejemplo.com", "debil", ErrClaveInvalida},
	}

	for _, caso := range casos {
		err := ValidarRegistro(caso.nombre, caso.correo, caso.clave)
		if !errors.Is(err, caso.espera) {
			t.Fatalf("ValidarRegistro() error = %v, espera %v", err, caso.espera)
		}
	}
}
