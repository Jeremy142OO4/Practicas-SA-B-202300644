package seguridad

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const LongitudClaveCorreo = 32

type GeneradorCorreo struct {
	clave []byte
}

func NuevoGeneradorCorreo(clave []byte) (*GeneradorCorreo, error) {
	if len(clave) != LongitudClaveCorreo {
		return nil, fmt.Errorf("la clave de correo debe tener %d bytes", LongitudClaveCorreo)
	}
	return &GeneradorCorreo{clave: append([]byte(nil), clave...)}, nil
}

func (generador *GeneradorCorreo) Generar(correo string) string {
	digest := hmac.New(sha256.New, generador.clave)
	_, _ = digest.Write([]byte(NormalizarCorreo(correo)))
	return hex.EncodeToString(digest.Sum(nil))
}

func NormalizarCorreo(correo string) string {
	return strings.ToLower(strings.TrimSpace(correo))
}

var _ IndiceCorreo = (*GeneradorCorreo)(nil)
