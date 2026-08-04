package token

type Generador interface {
	Generar(id string, rol string) (string, error)
}

type Validador interface {
	Validar(valor string) (*Datos, error)
}

type ValidadorRenovacion interface {
	ValidarRenovacion(valor string) (*Datos, error)
}
