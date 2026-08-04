package servicio

import (
	"fmt"

	"practica2/backend/internal/modules/users/dto"
	"practica2/backend/internal/modules/users/model"
	seguridad "practica2/backend/internal/shared/crypto"
)

func convertirUsuario(usuario modelo.Usuario, descifrador seguridad.Descifrador) (dto.Usuario, error) {
	nombre, err := descifrador.Descifrar(usuario.NombreCifrado)
	if err != nil {
		return dto.Usuario{}, fmt.Errorf("descifrar nombre: %w", err)
	}
	correo, err := descifrador.Descifrar(usuario.CorreoCifrado)
	if err != nil {
		return dto.Usuario{}, fmt.Errorf("descifrar correo: %w", err)
	}
	return dto.Usuario{
		ID: usuario.ID, Nombre: nombre, Correo: correo,
		Rol: usuario.Rol.Nombre, Estado: string(usuario.Estado),
		FechaCreacion: usuario.FechaCreacion, FechaActualiza: usuario.FechaActualiza,
	}, nil
}
