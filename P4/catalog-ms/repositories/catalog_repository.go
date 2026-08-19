package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"p4/catalog-ms/models"
)

var (
	ErrVideojuegoNoEncontrado   = errors.New("videojuego no encontrado")
	ErrCategoriaNoEncontrada    = errors.New("categoria no encontrada")
	ErrPlataformaNoEncontrada   = errors.New("plataforma no encontrada")
	ErrPresentacionNoEncontrada = errors.New(
		"videojuego en plataforma no encontrado",
	)
)

type RepositorioCatalogo interface {
	ListarVideojuegos() ([]models.Videojuego, error)

	ObtenerVideojuego(
		idVideojuego string,
	) (*models.Videojuego, error)

	CrearVideojuego(
		videojuego *models.Videojuego,
	) error

	ListarCategorias() ([]models.Categoria, error)

	ObtenerCategoria(
		idCategoria int,
	) (*models.Categoria, error)

	CrearCategoria(
		categoria *models.Categoria,
	) error

	ListarPlataformas() ([]models.Plataforma, error)

	ObtenerPlataforma(
		idPlataforma int,
	) (*models.Plataforma, error)

	CrearPlataforma(
		plataforma *models.Plataforma,
	) error

	CrearVideojuegoPlataforma(
		presentacion *models.VideojuegoPlataforma,
	) error

	ObtenerVideojuegoPlataforma(
		id string,
	) (*models.VideojuegoPlataforma, error)

	ListarPlataformasVideojuego(
		idVideojuego string,
	) ([]models.VideojuegoPlataforma, error)

	ActualizarStock(
		id string,
		stock int,
	) (*models.VideojuegoPlataforma, error)
}

type CatalogRepository struct {
	BD *pgxpool.Pool
}

func NuevoCatalogRepository(
	bd *pgxpool.Pool,
) *CatalogRepository {

	return &CatalogRepository{
		BD: bd,
	}
}

func (r *CatalogRepository) ListarVideojuegos() (
	[]models.Videojuego,
	error,
) {

	consulta := `
		SELECT
			id_videojuego,
			titulo,
			descripcion,
			id_categoria,
			desarrollador,
			fecha_lanzamiento,
			estado,
			created_at,
			updated_at
		FROM catalog.videojuegos
		ORDER BY titulo
	`

	filas, err := r.BD.Query(
		context.Background(),
		consulta,
	)

	if err != nil {
		return nil, err
	}

	defer filas.Close()

	videojuegos := []models.Videojuego{}

	for filas.Next() {

		var videojuego models.Videojuego

		err := filas.Scan(
			&videojuego.IDVideojuego,
			&videojuego.Titulo,
			&videojuego.Descripcion,
			&videojuego.IDCategoria,
			&videojuego.Desarrollador,
			&videojuego.FechaLanzamiento,
			&videojuego.Estado,
			&videojuego.CreatedAt,
			&videojuego.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		videojuegos = append(
			videojuegos,
			videojuego,
		)
	}

	return videojuegos, nil
}

func (r *CatalogRepository) ObtenerVideojuego(
	idVideojuego string,
) (*models.Videojuego, error) {

	consulta := `
		SELECT
			id_videojuego,
			titulo,
			descripcion,
			id_categoria,
			desarrollador,
			fecha_lanzamiento,
			estado,
			created_at,
			updated_at
		FROM catalog.videojuegos
		WHERE id_videojuego = $1
	`

	var videojuego models.Videojuego

	err := r.BD.QueryRow(
		context.Background(),
		consulta,
		idVideojuego,
	).Scan(
		&videojuego.IDVideojuego,
		&videojuego.Titulo,
		&videojuego.Descripcion,
		&videojuego.IDCategoria,
		&videojuego.Desarrollador,
		&videojuego.FechaLanzamiento,
		&videojuego.Estado,
		&videojuego.CreatedAt,
		&videojuego.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrVideojuegoNoEncontrado
	}

	return &videojuego, err
}

func (r *CatalogRepository) CrearVideojuego(
	videojuego *models.Videojuego,
) error {

	consulta := `
		INSERT INTO catalog.videojuegos (
			titulo,
			descripcion,
			id_categoria,
			desarrollador,
			fecha_lanzamiento,
			estado
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id_videojuego,
			created_at,
			updated_at
	`

	return r.BD.QueryRow(
		context.Background(),
		consulta,
		videojuego.Titulo,
		videojuego.Descripcion,
		videojuego.IDCategoria,
		videojuego.Desarrollador,
		videojuego.FechaLanzamiento,
		videojuego.Estado,
	).Scan(
		&videojuego.IDVideojuego,
		&videojuego.CreatedAt,
		&videojuego.UpdatedAt,
	)
}

func (r *CatalogRepository) ListarCategorias() (
	[]models.Categoria,
	error,
) {

	consulta := `
		SELECT
			id_categoria,
			nombre,
			descripcion
		FROM catalog.categorias
		ORDER BY nombre
	`

	filas, err := r.BD.Query(
		context.Background(),
		consulta,
	)

	if err != nil {
		return nil, err
	}

	defer filas.Close()

	categorias := []models.Categoria{}

	for filas.Next() {

		var categoria models.Categoria

		if err := filas.Scan(
			&categoria.IDCategoria,
			&categoria.Nombre,
			&categoria.Descripcion,
		); err != nil {

			return nil, err
		}

		categorias = append(
			categorias,
			categoria,
		)
	}

	return categorias, nil
}

func (r *CatalogRepository) ObtenerCategoria(
	idCategoria int,
) (*models.Categoria, error) {

	consulta := `
		SELECT
			id_categoria,
			nombre,
			descripcion
		FROM catalog.categorias
		WHERE id_categoria = $1
	`

	var categoria models.Categoria

	err := r.BD.QueryRow(
		context.Background(),
		consulta,
		idCategoria,
	).Scan(
		&categoria.IDCategoria,
		&categoria.Nombre,
		&categoria.Descripcion,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCategoriaNoEncontrada
	}

	return &categoria, err
}

func (r *CatalogRepository) CrearCategoria(
	categoria *models.Categoria,
) error {

	consulta := `
		INSERT INTO catalog.categorias (
			nombre,
			descripcion
		)
		VALUES ($1, $2)
		RETURNING id_categoria
	`

	return r.BD.QueryRow(
		context.Background(),
		consulta,
		categoria.Nombre,
		categoria.Descripcion,
	).Scan(
		&categoria.IDCategoria,
	)
}

func (r *CatalogRepository) ListarPlataformas() (
	[]models.Plataforma,
	error,
) {

	consulta := `
		SELECT
			id_plataforma,
			nombre,
			descripcion
		FROM catalog.plataformas
		ORDER BY nombre
	`

	filas, err := r.BD.Query(
		context.Background(),
		consulta,
	)

	if err != nil {
		return nil, err
	}

	defer filas.Close()

	plataformas := []models.Plataforma{}

	for filas.Next() {

		var plataforma models.Plataforma

		if err := filas.Scan(
			&plataforma.IDPlataforma,
			&plataforma.Nombre,
			&plataforma.Descripcion,
		); err != nil {

			return nil, err
		}

		plataformas = append(
			plataformas,
			plataforma,
		)
	}

	return plataformas, nil
}

func (r *CatalogRepository) ObtenerPlataforma(
	idPlataforma int,
) (*models.Plataforma, error) {

	consulta := `
		SELECT
			id_plataforma,
			nombre,
			descripcion
		FROM catalog.plataformas
		WHERE id_plataforma = $1
	`

	var plataforma models.Plataforma

	err := r.BD.QueryRow(
		context.Background(),
		consulta,
		idPlataforma,
	).Scan(
		&plataforma.IDPlataforma,
		&plataforma.Nombre,
		&plataforma.Descripcion,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPlataformaNoEncontrada
	}

	return &plataforma, err
}

func (r *CatalogRepository) CrearPlataforma(
	plataforma *models.Plataforma,
) error {

	consulta := `
		INSERT INTO catalog.plataformas (
			nombre,
			descripcion
		)
		VALUES ($1, $2)
		RETURNING id_plataforma
	`

	return r.BD.QueryRow(
		context.Background(),
		consulta,
		plataforma.Nombre,
		plataforma.Descripcion,
	).Scan(
		&plataforma.IDPlataforma,
	)
}

func (r *CatalogRepository) CrearVideojuegoPlataforma(
	presentacion *models.VideojuegoPlataforma,
) error {

	consulta := `
		INSERT INTO catalog.videojuego_plataforma (
			id_videojuego,
			id_plataforma,
			precio,
			stock,
			activo
		)
		VALUES ($1, $2, $3, $4, true)
		RETURNING
			id_videojuego_plataforma,
			activo
	`

	return r.BD.QueryRow(
		context.Background(),
		consulta,
		presentacion.IDVideojuego,
		presentacion.IDPlataforma,
		presentacion.Precio,
		presentacion.Stock,
	).Scan(
		&presentacion.IDVideojuegoPlataforma,
		&presentacion.Activo,
	)
}

func (r *CatalogRepository) ObtenerVideojuegoPlataforma(
	id string,
) (*models.VideojuegoPlataforma, error) {

	consulta := `
		SELECT
			vp.id_videojuego_plataforma,
			vp.id_videojuego,
			vp.id_plataforma,
			CAST(vp.precio AS DOUBLE PRECISION),
			vp.stock,
			vp.activo,
			p.id_plataforma,
			p.nombre,
			p.descripcion
		FROM catalog.videojuego_plataforma vp
		INNER JOIN catalog.plataformas p
			ON p.id_plataforma = vp.id_plataforma
		WHERE vp.id_videojuego_plataforma = $1
	`

	var presentacion models.VideojuegoPlataforma

	err := r.BD.QueryRow(
		context.Background(),
		consulta,
		id,
	).Scan(
		&presentacion.IDVideojuegoPlataforma,
		&presentacion.IDVideojuego,
		&presentacion.IDPlataforma,
		&presentacion.Precio,
		&presentacion.Stock,
		&presentacion.Activo,
		&presentacion.Plataforma.IDPlataforma,
		&presentacion.Plataforma.Nombre,
		&presentacion.Plataforma.Descripcion,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPresentacionNoEncontrada
	}

	return &presentacion, err
}

func (r *CatalogRepository) ListarPlataformasVideojuego(
	idVideojuego string,
) ([]models.VideojuegoPlataforma, error) {

	consulta := `
		SELECT
			vp.id_videojuego_plataforma,
			vp.id_videojuego,
			vp.id_plataforma,
			CAST(vp.precio AS DOUBLE PRECISION),
			vp.stock,
			vp.activo,
			p.id_plataforma,
			p.nombre,
			p.descripcion
		FROM catalog.videojuego_plataforma vp
		INNER JOIN catalog.plataformas p
			ON p.id_plataforma = vp.id_plataforma
		WHERE vp.id_videojuego = $1
		ORDER BY p.nombre
	`

	filas, err := r.BD.Query(
		context.Background(),
		consulta,
		idVideojuego,
	)

	if err != nil {
		return nil, err
	}

	defer filas.Close()

	resultado := []models.VideojuegoPlataforma{}

	for filas.Next() {

		var presentacion models.VideojuegoPlataforma

		if err := filas.Scan(
			&presentacion.IDVideojuegoPlataforma,
			&presentacion.IDVideojuego,
			&presentacion.IDPlataforma,
			&presentacion.Precio,
			&presentacion.Stock,
			&presentacion.Activo,
			&presentacion.Plataforma.IDPlataforma,
			&presentacion.Plataforma.Nombre,
			&presentacion.Plataforma.Descripcion,
		); err != nil {

			return nil, err
		}

		resultado = append(
			resultado,
			presentacion,
		)
	}

	return resultado, nil
}

func (r *CatalogRepository) ActualizarStock(
	id string,
	stock int,
) (*models.VideojuegoPlataforma, error) {

	consulta := `
		UPDATE catalog.videojuego_plataforma
		SET stock = $2
		WHERE id_videojuego_plataforma = $1
		RETURNING
			id_videojuego_plataforma,
			id_videojuego,
			id_plataforma,
			CAST(precio AS DOUBLE PRECISION),
			stock,
			activo
	`

	var presentacion models.VideojuegoPlataforma

	err := r.BD.QueryRow(
		context.Background(),
		consulta,
		id,
		stock,
	).Scan(
		&presentacion.IDVideojuegoPlataforma,
		&presentacion.IDVideojuego,
		&presentacion.IDPlataforma,
		&presentacion.Precio,
		&presentacion.Stock,
		&presentacion.Activo,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPresentacionNoEncontrada
	}

	return &presentacion, err
}
