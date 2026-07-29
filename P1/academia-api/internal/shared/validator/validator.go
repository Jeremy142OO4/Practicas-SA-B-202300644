package validator

import "strings"

func LimpiarTexto(value string) string {
	return strings.TrimSpace(value)
}

func ValidarTextoRequerido(
	value string,
	nombre string,
	longitudMaxima int,
) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return nombre + " es obligatorio"
	}

	if len([]rune(value)) > longitudMaxima {
		return nombre +
			" supera la longitud máxima permitida"
	}

	return ""
}
