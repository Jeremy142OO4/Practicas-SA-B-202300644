package validacion

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
)

var (
	ErrNombreInvalido = errors.New("el nombre debe tener entre 2 y 100 caracteres")
	ErrCorreoInvalido = errors.New("el correo no es valido")
	ErrClaveInvalida  = errors.New("la clave debe tener entre 8 y 72 caracteres, mayuscula, minuscula, numero y simbolo")
)

func ValidarRegistro(nombre string, correo string, clave string) error {
	nombre = strings.TrimSpace(nombre)
	if len([]rune(nombre)) < 2 || len([]rune(nombre)) > 100 {
		return ErrNombreInvalido
	}

	correo = strings.TrimSpace(correo)
	direccion, err := mail.ParseAddress(correo)
	if err != nil || direccion.Address != correo || len(correo) > 254 {
		return ErrCorreoInvalido
	}

	if !claveSegura(clave) {
		return ErrClaveInvalida
	}
	return nil
}

func claveSegura(clave string) bool {
	if len([]rune(clave)) < 8 || len([]rune(clave)) > 72 {
		return false
	}
	var mayuscula, minuscula, numero, simbolo bool
	for _, caracter := range clave {
		switch {
		case unicode.IsUpper(caracter):
			mayuscula = true
		case unicode.IsLower(caracter):
			minuscula = true
		case unicode.IsDigit(caracter):
			numero = true
		case unicode.IsPunct(caracter) || unicode.IsSymbol(caracter):
			simbolo = true
		}
	}
	return mayuscula && minuscula && numero && simbolo
}

func EsErrorRegistro(err error) bool {
	return errors.Is(err, ErrNombreInvalido) ||
		errors.Is(err, ErrCorreoInvalido) ||
		errors.Is(err, ErrClaveInvalida)
}
