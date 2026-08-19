package config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	IDUsuario string `json:"id_usuario"`
	Rol       string `json:"rol"`

	jwt.RegisteredClaims
}

type GestorJWT struct {
	Clave string
}

func NuevoGestorJWT(clave string) *GestorJWT {

	return &GestorJWT{
		Clave: clave,
	}
}

func (g *GestorJWT) GenerarToken(
	idUsuario string,
	rol string,
) (string, error) {

	claims := Claims{
		IDUsuario: idUsuario,
		Rol:       rol,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "auth-ms",

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(2 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(g.Clave),
	)
}

func (g *GestorJWT) ValidarToken(
	tokenString string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"método de token inválido",
				)
			}

			return []byte(g.Clave), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
