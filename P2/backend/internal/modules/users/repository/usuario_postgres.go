package repositorio

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"practica2/backend/internal/modules/users/model"
)

const consultaBuscarUsuario = `
SELECT
    u.id::text,
    u.name_encrypted,
    u.email_encrypted,
    u.email_hash,
    u.password_encrypted,
    r.id,
    r.name,
    r.description,
    u.status,
    u.created_at,
    u.updated_at
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.email_hash = $1`

const camposUsuario = `
    u.id::text,
    u.name_encrypted,
    u.email_encrypted,
    u.email_hash,
    u.password_encrypted,
    r.id,
    r.name,
    r.description,
    u.status,
    u.created_at,
    u.updated_at`

const consultaListarUsuarios = `SELECT ` + camposUsuario + `
FROM users u
JOIN roles r ON r.id = u.role_id
ORDER BY u.created_at ASC, u.id ASC`

const consultaBuscarPorID = `SELECT ` + camposUsuario + `
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.id = $1`

const consultaCambiarEstado = `
WITH usuario_actualizado AS (
    UPDATE users SET status = $2 WHERE id = $1 RETURNING *
)
SELECT ` + camposUsuario + `
FROM usuario_actualizado u
JOIN roles r ON r.id = u.role_id`

const consultaCrearUsuario = `
WITH rol_elegido AS (
    SELECT id
    FROM roles
    WHERE name = $5
), usuario_creado AS (
    INSERT INTO users (
        name_encrypted,
        email_encrypted,
        email_hash,
        password_encrypted,
        role_id,
        status
    )
    SELECT $1, $2, $3, $4, rol_elegido.id, $6
    FROM rol_elegido
    RETURNING *
)
SELECT
    u.id::text,
    u.name_encrypted,
    u.email_encrypted,
    u.email_hash,
    u.password_encrypted,
    r.id,
    r.name,
    r.description,
    u.status,
    u.created_at,
    u.updated_at
FROM usuario_creado u
JOIN roles r ON r.id = u.role_id`

type RepositorioPostgres struct {
	pool *pgxpool.Pool
}

func NuevoRepositorioPostgres(pool *pgxpool.Pool) (*RepositorioPostgres, error) {
	if pool == nil {
		return nil, errors.New("el pool de base de datos es obligatorio")
	}
	return &RepositorioPostgres{pool: pool}, nil
}

func (repositorio *RepositorioPostgres) BuscarPorIndice(ctx context.Context, indice string) (*modelo.Usuario, error) {
	fila := repositorio.pool.QueryRow(ctx, consultaBuscarUsuario, indice)
	usuario, err := leerUsuario(fila)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUsuarioNoExiste
	}
	if err != nil {
		return nil, fmt.Errorf("buscar usuario: %w", err)
	}
	return usuario, nil
}

func (repositorio *RepositorioPostgres) Crear(ctx context.Context, datos modelo.DatosUsuario) (*modelo.Usuario, error) {
	fila := repositorio.pool.QueryRow(
		ctx,
		consultaCrearUsuario,
		datos.NombreCifrado,
		datos.CorreoCifrado,
		datos.IndiceCorreo,
		datos.ClaveCifrada,
		datos.NombreRol,
		datos.Estado,
	)

	usuario, err := leerUsuario(fila)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRolNoExiste
	}
	if esCorreoDuplicado(err) {
		return nil, ErrCorreoDuplicado
	}
	if err != nil {
		return nil, fmt.Errorf("crear usuario: %w", err)
	}
	return usuario, nil
}

func (repositorio *RepositorioPostgres) Listar(ctx context.Context) ([]modelo.Usuario, error) {
	filas, err := repositorio.pool.Query(ctx, consultaListarUsuarios)
	if err != nil {
		return nil, fmt.Errorf("listar usuarios: %w", err)
	}
	defer filas.Close()
	usuarios := make([]modelo.Usuario, 0)
	for filas.Next() {
		usuario, err := leerUsuario(filas)
		if err != nil {
			return nil, fmt.Errorf("leer usuario: %w", err)
		}
		usuarios = append(usuarios, *usuario)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("recorrer usuarios: %w", err)
	}
	return usuarios, nil
}

func (repositorio *RepositorioPostgres) BuscarPorID(ctx context.Context, id string) (*modelo.Usuario, error) {
	usuario, err := leerUsuario(repositorio.pool.QueryRow(ctx, consultaBuscarPorID, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUsuarioNoExiste
	}
	if err != nil {
		return nil, fmt.Errorf("buscar usuario por id: %w", err)
	}
	return usuario, nil
}

func (repositorio *RepositorioPostgres) CambiarEstado(ctx context.Context, id string, estado modelo.EstadoUsuario) (*modelo.Usuario, error) {
	usuario, err := leerUsuario(repositorio.pool.QueryRow(ctx, consultaCambiarEstado, id, estado))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUsuarioNoExiste
	}
	if err != nil {
		return nil, fmt.Errorf("cambiar estado: %w", err)
	}
	return usuario, nil
}

type filaUsuario interface {
	Scan(destinos ...any) error
}

func leerUsuario(fila filaUsuario) (*modelo.Usuario, error) {
	var usuario modelo.Usuario
	var estado string
	err := fila.Scan(
		&usuario.ID,
		&usuario.NombreCifrado,
		&usuario.CorreoCifrado,
		&usuario.IndiceCorreo,
		&usuario.ClaveCifrada,
		&usuario.Rol.ID,
		&usuario.Rol.Nombre,
		&usuario.Rol.Descripcion,
		&estado,
		&usuario.FechaCreacion,
		&usuario.FechaActualiza,
	)
	if err != nil {
		return nil, err
	}

	usuario.Estado = modelo.EstadoUsuario(estado)
	if !usuario.Estado.EsValido() {
		return nil, fmt.Errorf("estado de usuario invalido: %s", estado)
	}
	return &usuario, nil
}

func esCorreoDuplicado(err error) bool {
	var errorPostgres *pgconn.PgError
	return errors.As(err, &errorPostgres) && errorPostgres.Code == "23505"
}

var _ RepositorioUsuario = (*RepositorioPostgres)(nil)
