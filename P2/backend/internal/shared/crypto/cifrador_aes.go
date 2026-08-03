package seguridad

import (
	"crypto/aes"
	"crypto/cipher"
	azar "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const LongitudClaveAES = 32

var ErrCifradoInvalido = errors.New("valor cifrado invalido")

type CifradorAES struct {
	aead cipher.AEAD
	azar io.Reader
}

func NuevoCifradorAES(clave []byte) (*CifradorAES, error) {
	return nuevoCifradorAES(clave, azar.Reader)
}

func nuevoCifradorAES(clave []byte, fuente io.Reader) (*CifradorAES, error) {
	if len(clave) != LongitudClaveAES {
		return nil, fmt.Errorf("la clave AES debe tener %d bytes", LongitudClaveAES)
	}
	if fuente == nil {
		return nil, errors.New("la fuente aleatoria es obligatoria")
	}

	bloque, err := aes.NewCipher(append([]byte(nil), clave...))
	if err != nil {
		return nil, fmt.Errorf("crear cifrador AES: %w", err)
	}
	aead, err := cipher.NewGCM(bloque)
	if err != nil {
		return nil, fmt.Errorf("crear cifrador GCM: %w", err)
	}
	return &CifradorAES{aead: aead, azar: fuente}, nil
}

func (servicio *CifradorAES) Cifrar(valor string) (string, error) {
	nonce := make([]byte, servicio.aead.NonceSize())
	if _, err := io.ReadFull(servicio.azar, nonce); err != nil {
		return "", fmt.Errorf("generar nonce: %w", err)
	}

	cifrado := servicio.aead.Seal(nil, nonce, []byte(valor), nil)
	contenido := make([]byte, 0, len(nonce)+len(cifrado))
	contenido = append(contenido, nonce...)
	contenido = append(contenido, cifrado...)
	return base64.StdEncoding.EncodeToString(contenido), nil
}

func (servicio *CifradorAES) Descifrar(valor string) (string, error) {
	contenido, err := base64.StdEncoding.DecodeString(valor)
	if err != nil {
		return "", ErrCifradoInvalido
	}
	minimo := servicio.aead.NonceSize() + servicio.aead.Overhead()
	if len(contenido) < minimo {
		return "", ErrCifradoInvalido
	}

	tamanoNonce := servicio.aead.NonceSize()
	texto, err := servicio.aead.Open(nil, contenido[:tamanoNonce], contenido[tamanoNonce:], nil)
	if err != nil {
		return "", ErrCifradoInvalido
	}
	return string(texto), nil
}

var _ Cifrador = (*CifradorAES)(nil)
var _ Descifrador = (*CifradorAES)(nil)
