package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalido = errors.New("token invalido")
	ErrTokenExpirado = errors.New("token expirado")
	ErrFueraDeGracia = errors.New("token fuera del tiempo de gracia")
)

type Gestor struct {
	secreto  []byte
	vigencia time.Duration
	gracia   time.Duration
	emisor   string
	reloj    func() time.Time
}

func NuevoGestor(secreto string, vigencia time.Duration, gracia time.Duration, emisor string) (*Gestor, error) {
	if len(secreto) < 32 {
		return nil, errors.New("el secreto JWT debe tener al menos 32 caracteres")
	}
	if vigencia <= 0 || gracia <= 0 {
		return nil, errors.New("la vigencia y la gracia deben ser positivas")
	}
	if emisor == "" {
		return nil, errors.New("el emisor JWT es obligatorio")
	}
	return &Gestor{
		secreto: []byte(secreto), vigencia: vigencia, gracia: gracia,
		emisor: emisor, reloj: time.Now,
	}, nil
}

func (gestor *Gestor) Generar(id string, rol string) (string, error) {
	if id == "" || rol == "" {
		return "", errors.New("el usuario y el rol son obligatorios")
	}
	ahora := gestor.reloj()
	claims := Claims{
		Rol: rol,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: id, Issuer: gestor.emisor,
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(ahora.Add(gestor.vigencia)),
		},
	}
	valor, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(gestor.secreto)
	if err != nil {
		return "", fmt.Errorf("firmar token: %w", err)
	}
	return valor, nil
}

func (gestor *Gestor) Validar(valor string) (*Datos, error) {
	claims, err := gestor.leer(valor, false)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrTokenExpirado
	}
	if err != nil {
		return nil, ErrTokenInvalido
	}
	datos, err := convertir(claims)
	if err == nil {
		datos.RenovableHasta = claims.ExpiresAt.Time.Add(gestor.gracia).Unix()
	}
	return datos, err
}

func (gestor *Gestor) ValidarRenovacion(valor string) (*Datos, error) {
	claims, err := gestor.leer(valor, true)
	if err != nil {
		return nil, ErrTokenInvalido
	}
	if claims.ExpiresAt == nil || claims.Subject == "" || claims.Rol == "" || claims.Issuer != gestor.emisor {
		return nil, ErrTokenInvalido
	}
	ahora := gestor.reloj()
	if ahora.Before(claims.ExpiresAt.Time) {
		return nil, ErrTokenInvalido
	}
	if ahora.After(claims.ExpiresAt.Time.Add(gestor.gracia)) {
		return nil, ErrFueraDeGracia
	}
	datos, err := convertir(claims)
	if err == nil {
		datos.RenovableHasta = claims.ExpiresAt.Time.Add(gestor.gracia).Unix()
	}
	return datos, err
}

func (gestor *Gestor) leer(valor string, omitirValidacion bool) (*Claims, error) {
	claims := &Claims{}
	opciones := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(gestor.reloj),
	}
	if omitirValidacion {
		opciones = append(opciones, jwt.WithoutClaimsValidation())
	} else {
		opciones = append(opciones, jwt.WithIssuer(gestor.emisor), jwt.WithExpirationRequired())
	}
	tokenFirmado, err := jwt.ParseWithClaims(valor, claims, func(token *jwt.Token) (any, error) {
		return gestor.secreto, nil
	}, opciones...)
	if err != nil || tokenFirmado == nil || !tokenFirmado.Valid {
		return nil, err
	}
	return claims, nil
}

func convertir(claims *Claims) (*Datos, error) {
	if claims.Subject == "" || claims.Rol == "" || claims.ExpiresAt == nil {
		return nil, ErrTokenInvalido
	}
	return &Datos{ID: claims.Subject, Rol: claims.Rol, ExpiraEn: claims.ExpiresAt.Unix()}, nil
}

var _ Generador = (*Gestor)(nil)
var _ Validador = (*Gestor)(nil)
var _ ValidadorRenovacion = (*Gestor)(nil)
