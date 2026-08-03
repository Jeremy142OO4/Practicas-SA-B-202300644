package configuracion

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configuracion struct {
	Aplicacion Aplicacion
	BaseDatos  BaseDatos
	JWT        JWT
	Cookie     Cookie
	Seguridad  Seguridad
}

type Aplicacion struct {
	Entorno     string
	Puerto      int
	URLFrontend string
}

type BaseDatos struct {
	URL string
}

type JWT struct {
	Secreto           string
	MinutosExpiracion int
	MinutosGracia     int
	Emisor            string
}

type Cookie struct {
	Segura   bool
	SameSite string
}

type Seguridad struct {
	ClaveAES    []byte
	ClaveCorreo []byte
}

func Cargar() (Configuracion, error) {
	_ = godotenv.Load()

	puerto, err := enteroPositivo("APP_PORT")
	if err != nil {
		return Configuracion{}, err
	}
	expiracion, err := enteroPositivo("JWT_EXPIRATION_MINUTES")
	if err != nil {
		return Configuracion{}, err
	}
	gracia, err := enteroPositivo("JWT_RENEWAL_GRACE_MINUTES")
	if err != nil {
		return Configuracion{}, err
	}
	segura, err := valorBooleano("COOKIE_SECURE")
	if err != nil {
		return Configuracion{}, err
	}
	claveAES, err := claveBase64("AES_KEY", 32)
	if err != nil {
		return Configuracion{}, err
	}
	claveCorreo, err := claveBase64("EMAIL_INDEX_KEY", 32)
	if err != nil {
		return Configuracion{}, err
	}

	config := Configuracion{
		Aplicacion: Aplicacion{
			Entorno:     valorPredeterminado("APP_ENV", "development"),
			Puerto:      puerto,
			URLFrontend: strings.TrimSpace(os.Getenv("FRONTEND_URL")),
		},
		BaseDatos: BaseDatos{URL: strings.TrimSpace(os.Getenv("DATABASE_URL"))},
		JWT: JWT{
			Secreto:           strings.TrimSpace(os.Getenv("JWT_SECRET")),
			MinutosExpiracion: expiracion,
			MinutosGracia:     gracia,
			Emisor:            strings.TrimSpace(os.Getenv("JWT_ISSUER")),
		},
		Cookie: Cookie{
			Segura:   segura,
			SameSite: strings.ToLower(strings.TrimSpace(os.Getenv("COOKIE_SAME_SITE"))),
		},
		Seguridad: Seguridad{ClaveAES: claveAES, ClaveCorreo: claveCorreo},
	}

	if err := config.validar(); err != nil {
		return Configuracion{}, err
	}
	return config, nil
}

func (config Configuracion) validar() error {
	requeridas := map[string]string{
		"FRONTEND_URL": config.Aplicacion.URLFrontend,
		"DATABASE_URL": config.BaseDatos.URL,
		"JWT_SECRET":   config.JWT.Secreto,
		"JWT_ISSUER":   config.JWT.Emisor,
	}
	for nombre, valor := range requeridas {
		if valor == "" {
			return fmt.Errorf("la variable %s es obligatoria", nombre)
		}
	}

	if len(config.JWT.Secreto) < 32 {
		return errors.New("JWT_SECRET debe tener al menos 32 caracteres")
	}
	valoresSameSite := map[string]bool{"lax": true, "strict": true, "none": true}
	if !valoresSameSite[config.Cookie.SameSite] {
		return errors.New("COOKIE_SAME_SITE debe ser Lax, Strict o None")
	}
	if config.Cookie.SameSite == "none" && !config.Cookie.Segura {
		return errors.New("COOKIE_SECURE debe ser true cuando COOKIE_SAME_SITE es None")
	}
	return nil
}

func enteroPositivo(nombre string) (int, error) {
	valor := strings.TrimSpace(os.Getenv(nombre))
	numero, err := strconv.Atoi(valor)
	if err != nil || numero <= 0 {
		return 0, fmt.Errorf("la variable %s debe ser un entero positivo", nombre)
	}
	return numero, nil
}

func valorBooleano(nombre string) (bool, error) {
	valor, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(nombre)))
	if err != nil {
		return false, fmt.Errorf("la variable %s debe ser true o false", nombre)
	}
	return valor, nil
}

func claveBase64(nombre string, longitud int) ([]byte, error) {
	clave, err := base64.StdEncoding.DecodeString(strings.TrimSpace(os.Getenv(nombre)))
	if err != nil || len(clave) != longitud {
		return nil, fmt.Errorf("la variable %s debe ser una clave Base64 de %d bytes", nombre, longitud)
	}
	return clave, nil
}

func valorPredeterminado(nombre string, predeterminado string) string {
	valor := strings.TrimSpace(os.Getenv(nombre))
	if valor == "" {
		return predeterminado
	}
	return valor
}
