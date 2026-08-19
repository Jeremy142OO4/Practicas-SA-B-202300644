package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"p4/auth-ms/models"
)

var (
	ErrUsuarioNoEncontrado = errors.New(
		"usuario no encontrado",
	)

	ErrCorreoRegistrado = errors.New(
		"el correo ya está registrado",
	)
)

type RepositorioUsuario interface {
	CrearUsuario(
		usuario *models.Usuario,
	) error

	BuscarPorHashCorreo(
		hashCorreo string,
	) (*models.Usuario, error)
}

type UsuarioRepository struct {
	BD *pgxpool.Pool
}

func NuevoUsuarioRepository(
	bd *pgxpool.Pool,
) *UsuarioRepository {

	return &UsuarioRepository{
		BD: bd,
	}
}

func (r *UsuarioRepository) CrearUsuario(
	usuario *models.Usuario,
) error {

	consulta := `
		INSERT INTO auth.usuarios (
			name_encrypted,
			email_encrypted,
			email_hash,
			password_encrypted,
			role_id,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id_usuario,
			created_at,
			updated_at
	`

	err := r.BD.QueryRow(
		context.Background(),
		consulta,

		usuario.NameEncrypted,
		usuario.EmailEncrypted,
		usuario.EmailHash,
		usuario.PasswordEncrypted,
		usuario.RoleID,
		usuario.Status,
	).Scan(
		&usuario.IDUsuario,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if esDatoDuplicado(err) {
		return ErrCorreoRegistrado
	}

	return err
}

func (r *UsuarioRepository) BuscarPorHashCorreo(
	hashCorreo string,
) (*models.Usuario, error) {

	consulta := `
		SELECT
			u.id_usuario,
			u.name_encrypted,
			u.email_encrypted,
			u.email_hash,
			u.password_encrypted,
			u.role_id,
			u.status,
			u.created_at,
			u.updated_at,
			r.id_rol,
			r.nombre,
			r.descripcion
		FROM auth.usuarios u
		INNER JOIN auth.roles r
			ON r.id_rol = u.role_id
		WHERE u.email_hash = $1
	`

	var usuario models.Usuario

	err := r.BD.QueryRow(
		context.Background(),
		consulta,
		hashCorreo,
	).Scan(
		&usuario.IDUsuario,
		&usuario.NameEncrypted,
		&usuario.EmailEncrypted,
		&usuario.EmailHash,
		&usuario.PasswordEncrypted,
		&usuario.RoleID,
		&usuario.Status,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,

		&usuario.Rol.IDRol,
		&usuario.Rol.Nombre,
		&usuario.Rol.Descripcion,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUsuarioNoEncontrado
	}

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

func esDatoDuplicado(err error) bool {

	if err == nil {
		return false
	}

	var errorPostgres *pgconn.PgError

	if errors.As(err, &errorPostgres) {
		return errorPostgres.Code == "23505"
	}

	return false
}
