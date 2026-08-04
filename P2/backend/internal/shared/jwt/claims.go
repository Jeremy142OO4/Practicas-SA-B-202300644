package token

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Rol string `json:"role"`
	jwt.RegisteredClaims
}

type Datos struct {
	ID             string
	Rol            string
	ExpiraEn       int64
	RenovableHasta int64
}
