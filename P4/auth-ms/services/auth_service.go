package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"p4/auth-ms/config"
	"p4/auth-ms/models"
	"p4/auth-ms/repositories"
)

var (
	ErrDatosInvalidos = errors.New(
		"nombre, correo y contraseña son obligatorios",
	)

	ErrCredencialesInvalidas = errors.New(
		"correo o contraseña incorrectos",
	)

	ErrUsuarioInactivo = errors.New(
		"el usuario está inactivo",
	)
)

type AuthService struct {
	Repositorio  repositories.RepositorioUsuario
	GestorJWT    *config.GestorJWT
	ClaveCifrado string
}

func NuevoAuthService(
	repositorio repositories.RepositorioUsuario,
	gestorJWT *config.GestorJWT,
	claveCifrado string,
) *AuthService {

	return &AuthService{
		Repositorio:  repositorio,
		GestorJWT:    gestorJWT,
		ClaveCifrado: claveCifrado,
	}
}

func (s *AuthService) RegistrarUsuario(
	datos models.RegistroRequest,
) (*models.UsuarioResponse, error) {

	datos.Nombre = strings.TrimSpace(
		datos.Nombre,
	)

	datos.Correo = strings.ToLower(
		strings.TrimSpace(datos.Correo),
	)

	if datos.Nombre == "" ||
		datos.Correo == "" ||
		datos.Contrasena == "" {

		return nil, ErrDatosInvalidos
	}

	hashCorreo := generarHashCorreo(
		datos.Correo,
		s.ClaveCifrado,
	)

	usuarioExistente, err :=
		s.Repositorio.BuscarPorHashCorreo(
			hashCorreo,
		)

	if err == nil && usuarioExistente != nil {
		return nil, repositories.ErrCorreoRegistrado
	}

	if err != nil &&
		!errors.Is(
			err,
			repositories.ErrUsuarioNoEncontrado,
		) {

		return nil, err
	}

	nombreCifrado, err := cifrarTexto(
		datos.Nombre,
		s.ClaveCifrado,
	)

	if err != nil {
		return nil, err
	}

	correoCifrado, err := cifrarTexto(
		datos.Correo,
		s.ClaveCifrado,
	)

	if err != nil {
		return nil, err
	}

	hashContrasena, err := bcrypt.GenerateFromPassword(
		[]byte(datos.Contrasena),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	usuario := &models.Usuario{
		NameEncrypted: nombreCifrado,

		EmailEncrypted: correoCifrado,

		EmailHash: hashCorreo,

		PasswordEncrypted: string(
			hashContrasena,
		),

		// Rol CLIENTE
		RoleID: 2,

		Status: "ACTIVO",
	}

	err = s.Repositorio.CrearUsuario(
		usuario,
	)

	if err != nil {
		return nil, err
	}

	respuesta := &models.UsuarioResponse{
		IDUsuario: usuario.IDUsuario,
		Nombre:    datos.Nombre,
		Correo:    datos.Correo,
		Rol:       "CLIENTE",
		Status:    usuario.Status,
	}

	return respuesta, nil
}

func (s *AuthService) IniciarSesion(
	datos models.LoginRequest,
) (*models.LoginResponse, error) {

	datos.Correo = strings.ToLower(
		strings.TrimSpace(datos.Correo),
	)

	if datos.Correo == "" ||
		datos.Contrasena == "" {

		return nil, ErrCredencialesInvalidas
	}

	hashCorreo := generarHashCorreo(
		datos.Correo,
		s.ClaveCifrado,
	)

	usuario, err :=
		s.Repositorio.BuscarPorHashCorreo(
			hashCorreo,
		)

	if err != nil {

		if errors.Is(
			err,
			repositories.ErrUsuarioNoEncontrado,
		) {

			return nil, ErrCredencialesInvalidas
		}

		return nil, err
	}

	if usuario.Status != "ACTIVO" {
		return nil, ErrUsuarioInactivo
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(usuario.PasswordEncrypted),
		[]byte(datos.Contrasena),
	)

	if err != nil {
		return nil, ErrCredencialesInvalidas
	}

	nombre, err := descifrarTexto(
		usuario.NameEncrypted,
		s.ClaveCifrado,
	)

	if err != nil {
		return nil, err
	}

	correo, err := descifrarTexto(
		usuario.EmailEncrypted,
		s.ClaveCifrado,
	)

	if err != nil {
		return nil, err
	}

	token, err := s.GestorJWT.GenerarToken(
		usuario.IDUsuario,
		usuario.Rol.Nombre,
	)

	if err != nil {
		return nil, err
	}

	respuesta := &models.LoginResponse{

		Token: token,

		Usuario: models.UsuarioResponse{
			IDUsuario: usuario.IDUsuario,
			Nombre:    nombre,
			Correo:    correo,
			Rol:       usuario.Rol.Nombre,
			Status:    usuario.Status,
		},
	}

	return respuesta, nil
}

func generarHashCorreo(
	correo string,
	clave string,
) string {

	h := hmac.New(
		sha256.New,
		[]byte(clave),
	)

	h.Write(
		[]byte(
			strings.ToLower(
				strings.TrimSpace(correo),
			),
		),
	)

	return hex.EncodeToString(
		h.Sum(nil),
	)
}

func cifrarTexto(
	texto string,
	clave string,
) (string, error) {

	claveAES := sha256.Sum256(
		[]byte(clave),
	)

	bloque, err := aes.NewCipher(
		claveAES[:],
	)

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(
		bloque,
	)

	if err != nil {
		return "", err
	}

	nonce := make(
		[]byte,
		gcm.NonceSize(),
	)

	if _, err := io.ReadFull(
		rand.Reader,
		nonce,
	); err != nil {

		return "", err
	}

	datos := gcm.Seal(
		nonce,
		nonce,
		[]byte(texto),
		nil,
	)

	return base64.StdEncoding.EncodeToString(
		datos,
	), nil
}

func descifrarTexto(
	textoCifrado string,
	clave string,
) (string, error) {

	claveAES := sha256.Sum256(
		[]byte(clave),
	)

	bloque, err := aes.NewCipher(
		claveAES[:],
	)

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(
		bloque,
	)

	if err != nil {
		return "", err
	}

	datos, err :=
		base64.StdEncoding.DecodeString(
			textoCifrado,
		)

	if err != nil {
		return "", err
	}

	tamanoNonce := gcm.NonceSize()

	if len(datos) < tamanoNonce {
		return "", errors.New(
			"texto cifrado inválido",
		)
	}

	nonce := datos[:tamanoNonce]

	contenido := datos[tamanoNonce:]

	texto, err := gcm.Open(
		nil,
		nonce,
		contenido,
		nil,
	)

	if err != nil {
		return "", err
	}

	return string(texto), nil
}
