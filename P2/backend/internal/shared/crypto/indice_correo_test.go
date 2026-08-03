package seguridad

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestIndiceNormalizaCorreo(t *testing.T) {
	generador := crearGenerador(t, bytes.Repeat([]byte{0x76}, LongitudClaveCorreo))
	primero := generador.Generar(" Usuario@Ejemplo.COM ")
	segundo := generador.Generar("usuario@ejemplo.com")
	if primero != segundo || len(primero) != 64 {
		t.Fatalf("indices no coinciden: %q != %q", primero, segundo)
	}
}

func TestIndiceDistingueCorreos(t *testing.T) {
	generador := crearGenerador(t, bytes.Repeat([]byte{0x87}, LongitudClaveCorreo))
	if generador.Generar("uno@ejemplo.com") == generador.Generar("dos@ejemplo.com") {
		t.Fatal("Generar() devolvio el mismo indice")
	}
}

func TestIndiceValidaClave(t *testing.T) {
	if _, err := NuevoGeneradorCorreo([]byte("corta")); err == nil {
		t.Fatal("NuevoGeneradorCorreo() debio rechazar la clave")
	}
}

func TestIndiceCoincideConSemillas(t *testing.T) {
	clave, _ := base64.StdEncoding.DecodeString("uk+bcsNpQJeuV6mmOIWQoqAPQ7icnYZt4sKul7d2Qsk=")
	generador := crearGenerador(t, clave)
	casos := map[string]string{
		"admin@practica.local":     "abf0788e2783e4a0caf57c4725c4488dbba89ff1891cc5823d8f41247ec667e6",
		" CLIENTE@PRACTICA.LOCAL ": "3c42f43ca98a926d5da41c3a69693f9b587139cbfa69371915b1870031f34311",
	}
	for correo, espera := range casos {
		if indice := generador.Generar(correo); indice != espera {
			t.Fatalf("Generar(%q) = %q", correo, indice)
		}
	}
}

func crearGenerador(t *testing.T, clave []byte) *GeneradorCorreo {
	t.Helper()
	generador, err := NuevoGeneradorCorreo(clave)
	if err != nil {
		t.Fatalf("NuevoGeneradorCorreo() error = %v", err)
	}
	return generador
}
