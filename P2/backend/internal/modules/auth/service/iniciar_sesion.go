package servicio

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"

	"practica2/backend/internal/modules/auth/dto"
	"practica2/backend/internal/modules/users/model"
	repositorio "practica2/backend/internal/modules/users/repository"
	seguridad "practica2/backend/internal/shared/crypto"
	token "practica2/backend/internal/shared/jwt"
)

var (
	ErrCredenciales    = errors.New("credenciales incorrectas")
	ErrUsuarioInactivo = errors.New("usuario inactivo")
)

type InicioSesiones interface {
	Iniciar(ctx context.Context, solicitud dto.SolicitudLogin) (*ResultadoLogin, error)
}

type ResultadoLogin struct {
	Token   string
	Usuario dto.UsuarioSesion
}

type ServicioLogin struct {
	lector      repositorio.LectorUsuario
	descifrador seguridad.Descifrador
	generador   token.Generador
	indice      seguridad.IndiceCorreo
}

func NuevoServicioLogin(
	lector repositorio.LectorUsuario,
	descifrador seguridad.Descifrador,
	generador token.Generador,
	indice seguridad.IndiceCorreo,
) (*ServicioLogin, error) {
	if lector == nil || descifrador == nil || generador == nil || indice == nil {
		return nil, errors.New("las dependencias de login son obligatorias")
	}
	return &ServicioLogin{lector: lector, descifrador: descifrador, generador: generador, indice: indice}, nil
}

func (servicio *ServicioLogin) Iniciar(
	ctx context.Context,
	solicitud dto.SolicitudLogin,
) (*ResultadoLogin, error) {
	correo := seguridad.NormalizarCorreo(solicitud.Correo)
	if correo == "" || solicitud.Clave == "" {
		return nil, ErrCredenciales
	}

	usuario, err := servicio.lector.BuscarPorIndice(ctx, servicio.indice.Generar(correo))
	if errors.Is(err, repositorio.ErrUsuarioNoExiste) {
		return nil, ErrCredenciales
	}
	if err != nil {
		return nil, fmt.Errorf("buscar credenciales: %w", err)
	}
	if usuario.Estado != modelo.EstadoActivo {
		return nil, ErrUsuarioInactivo
	}

	clave, err := servicio.descifrador.Descifrar(usuario.ClaveCifrada)
	if err != nil {
		return nil, fmt.Errorf("descifrar credenciales: %w", err)
	}
	if subtle.ConstantTimeCompare([]byte(clave), []byte(solicitud.Clave)) != 1 {
		return nil, ErrCredenciales
	}

	nombre, err := servicio.descifrador.Descifrar(usuario.NombreCifrado)
	if err != nil {
		return nil, fmt.Errorf("descifrar nombre: %w", err)
	}
	correoGuardado, err := servicio.descifrador.Descifrar(usuario.CorreoCifrado)
	if err != nil {
		return nil, fmt.Errorf("descifrar correo: %w", err)
	}
	valorToken, err := servicio.generador.Generar(usuario.ID, usuario.Rol.Nombre)
	if err != nil {
		return nil, fmt.Errorf("generar token: %w", err)
	}

	return &ResultadoLogin{
		Token: valorToken,
		Usuario: dto.UsuarioSesion{
			ID: usuario.ID, Nombre: nombre, Correo: correoGuardado,
			Rol: usuario.Rol.Nombre, Estado: string(usuario.Estado),
		},
	}, nil
}

var _ InicioSesiones = (*ServicioLogin)(nil)
