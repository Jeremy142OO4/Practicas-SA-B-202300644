package seguridad

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
)

func TestCifrarYDescifrar(t *testing.T) {
	servicio := crearCifrador(t, bytes.Repeat([]byte{0x21}, LongitudClaveAES))
	cifrado, err := servicio.Cifrar("valor sensible")
	if err != nil {
		t.Fatalf("Cifrar() devolvio error: %v", err)
	}
	if cifrado == "valor sensible" {
		t.Fatal("Cifrar() devolvio texto plano")
	}
	texto, err := servicio.Descifrar(cifrado)
	if err != nil {
		t.Fatalf("Descifrar() devolvio error: %v", err)
	}
	if texto != "valor sensible" {
		t.Fatalf("Descifrar() = %q", texto)
	}
}

func TestCifradoUsaNonceUnico(t *testing.T) {
	servicio := crearCifrador(t, bytes.Repeat([]byte{0x32}, LongitudClaveAES))
	primero, _ := servicio.Cifrar("mismo valor")
	segundo, _ := servicio.Cifrar("mismo valor")
	if primero == segundo {
		t.Fatal("Cifrar() reutilizo el nonce")
	}
}

func TestCifradoRechazaCambios(t *testing.T) {
	servicio := crearCifrador(t, bytes.Repeat([]byte{0x43}, LongitudClaveAES))
	cifrado, _ := servicio.Cifrar("no modificar")
	contenido, _ := base64.StdEncoding.DecodeString(cifrado)
	contenido[len(contenido)-1] ^= 0xff
	_, err := servicio.Descifrar(base64.StdEncoding.EncodeToString(contenido))
	if !errors.Is(err, ErrCifradoInvalido) {
		t.Fatalf("Descifrar() error = %v", err)
	}
}

func TestCifradoRechazaOtraClave(t *testing.T) {
	primero := crearCifrador(t, bytes.Repeat([]byte{0x54}, LongitudClaveAES))
	segundo := crearCifrador(t, bytes.Repeat([]byte{0x65}, LongitudClaveAES))
	cifrado, _ := primero.Cifrar("protegido")
	_, err := segundo.Descifrar(cifrado)
	if !errors.Is(err, ErrCifradoInvalido) {
		t.Fatalf("Descifrar() error = %v", err)
	}
}

func TestCifradoValidaClave(t *testing.T) {
	if _, err := NuevoCifradorAES([]byte("corta")); err == nil {
		t.Fatal("NuevoCifradorAES() debio rechazar la clave")
	}
}

func TestCifradoLeeSemillas(t *testing.T) {
	clave, _ := base64.StdEncoding.DecodeString("j2xztqsjTNGsofBPhs8bLn616i7vlOFQt3VUj21KbCc=")
	servicio := crearCifrador(t, clave)
	casos := []struct {
		nombre  string
		cifrado string
		espera  string
	}{
		{"correo admin", "I0Y9WZVwurC1xr0F3IpkY7ykbr7PUEnJaXdNtdwedOi++QZEefR6L+ZmCJh/RHhL", "admin@practica.local"},
		{"correo cliente", "ULa0xtTtRHOWOFwi1CAws0rd398d9JhAQsO7m3FJoOxkFiTUlq+XpXPlp4bVc9sdroA=", "cliente@practica.local"},
		{"clave admin", "EfpD+gHcLOn+ZxUGqJAVMBUmrXg4iUmqGRfq0tK78Zzl4unb5w==", "Admin123!"},
		{"clave cliente", "Zw+NH7AOHqWGEQpp8iLI6kVMWgQBwnY5EuI0N1RiJzuLY2eRQY/f", "Cliente123!"},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			texto, err := servicio.Descifrar(caso.cifrado)
			if err != nil || texto != caso.espera {
				t.Fatalf("Descifrar() = %q, error = %v", texto, err)
			}
		})
	}
}

func crearCifrador(t *testing.T, clave []byte) *CifradorAES {
	t.Helper()
	servicio, err := NuevoCifradorAES(clave)
	if err != nil {
		t.Fatalf("NuevoCifradorAES() error = %v", err)
	}
	return servicio
}
