package servicio

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"practica2/backend/internal/modules/auth/dto"
	"practica2/backend/internal/modules/users/model"
	repositorio "practica2/backend/internal/modules/users/repository"
	seguridad "practica2/backend/internal/shared/crypto"
	validacion "practica2/backend/internal/shared/validator"
)

type RegistroUsuarios interface {
	Registrar(ctx context.Context, solicitud dto.SolicitudRegistro) (*dto.UsuarioRegistrado, error)
}

type ServicioRegistro struct {
	lector    repositorio.LectorUsuario
	escritor  repositorio.EscritorUsuario
	cifrador  seguridad.Cifrador
	generador seguridad.IndiceCorreo
}

func NuevoServicioRegistro(
	lector repositorio.LectorUsuario,
	escritor repositorio.EscritorUsuario,
	cifrador seguridad.Cifrador,
	generador seguridad.IndiceCorreo,
) (*ServicioRegistro, error) {
	if lector == nil || escritor == nil || cifrador == nil || generador == nil {
		return nil, errors.New("las dependencias de registro son obligatorias")
	}
	return &ServicioRegistro{
		lector: lector, escritor: escritor, cifrador: cifrador, generador: generador,
	}, nil
}

func (servicio *ServicioRegistro) Registrar(
	ctx context.Context,
	solicitud dto.SolicitudRegistro,
) (*dto.UsuarioRegistrado, error) {
	if err := validacion.ValidarRegistro(solicitud.Nombre, solicitud.Correo, solicitud.Clave); err != nil {
		return nil, err
	}

	nombre := strings.TrimSpace(solicitud.Nombre)
	correo := seguridad.NormalizarCorreo(solicitud.Correo)
	indice := servicio.generador.Generar(correo)

	_, err := servicio.lector.BuscarPorIndice(ctx, indice)
	switch {
	case err == nil:
		return nil, repositorio.ErrCorreoDuplicado
	case !errors.Is(err, repositorio.ErrUsuarioNoExiste):
		return nil, fmt.Errorf("comprobar correo: %w", err)
	}

	nombreCifrado, err := servicio.cifrador.Cifrar(nombre)
	if err != nil {
		return nil, fmt.Errorf("cifrar nombre: %w", err)
	}
	correoCifrado, err := servicio.cifrador.Cifrar(correo)
	if err != nil {
		return nil, fmt.Errorf("cifrar correo: %w", err)
	}
	claveCifrada, err := servicio.cifrador.Cifrar(solicitud.Clave)
	if err != nil {
		return nil, fmt.Errorf("cifrar clave: %w", err)
	}

	usuario, err := servicio.escritor.Crear(ctx, modelo.DatosUsuario{
		NombreCifrado: nombreCifrado,
		CorreoCifrado: correoCifrado,
		IndiceCorreo:  indice,
		ClaveCifrada:  claveCifrada,
		NombreRol:     "CLIENTE",
		Estado:        modelo.EstadoActivo,
	})
	if err != nil {
		return nil, err
	}

	return &dto.UsuarioRegistrado{
		ID: usuario.ID, Nombre: nombre, Correo: correo,
		Rol: usuario.Rol.Nombre, Estado: string(usuario.Estado),
	}, nil
}

var _ RegistroUsuarios = (*ServicioRegistro)(nil)
