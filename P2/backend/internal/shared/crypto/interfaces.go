package seguridad

type Cifrador interface {
	Cifrar(valor string) (string, error)
}

type Descifrador interface {
	Descifrar(valor string) (string, error)
}

type IndiceCorreo interface {
	Generar(correo string) string
}
