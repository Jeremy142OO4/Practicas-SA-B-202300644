package services

import (
	"errors"
	"strings"
	"time"

	"p4/catalog-ms/models"
	"p4/catalog-ms/repositories"
)

var (
	ErrDatosInvalidos = errors.New(
		"los datos enviados son inválidos",
	)

	ErrStockInvalido = errors.New(
		"el stock no puede ser negativo",
	)

	ErrPrecioInvalido = errors.New(
		"el precio debe ser mayor o igual a cero",
	)
)

type CatalogService struct {
	Repositorio repositories.RepositorioCatalogo
}

func NuevoCatalogService(
	repositorio repositories.RepositorioCatalogo,
) *CatalogService {

	return &CatalogService{
		Repositorio: repositorio,
	}
}

func (s *CatalogService) ListarVideojuegos() (
	[]models.Videojuego,
	error,
) {

	return s.Repositorio.ListarVideojuegos()
}

func (s *CatalogService) ObtenerVideojuego(
	idVideojuego string,
) (*models.Videojuego, error) {

	return s.Repositorio.ObtenerVideojuego(
		idVideojuego,
	)
}

func (s *CatalogService) CrearVideojuego(
	datos models.CrearVideojuegoInput,
) (*models.Videojuego, error) {

	datos.Titulo = strings.TrimSpace(
		datos.Titulo,
	)

	datos.Desarrollador = strings.TrimSpace(
		datos.Desarrollador,
	)

	if datos.Titulo == "" ||
		datos.IDCategoria <= 0 ||
		datos.Desarrollador == "" ||
		datos.FechaLanzamiento == "" {

		return nil, ErrDatosInvalidos
	}

	_, err := s.Repositorio.ObtenerCategoria(
		datos.IDCategoria,
	)

	if err != nil {
		return nil, err
	}

	fecha, err := time.Parse(
		"2006-01-02",
		datos.FechaLanzamiento,
	)

	if err != nil {
		return nil, errors.New(
			"fecha de lanzamiento inválida",
		)
	}

	videojuego := &models.Videojuego{
		Titulo:           datos.Titulo,
		Descripcion:      datos.Descripcion,
		IDCategoria:      datos.IDCategoria,
		Desarrollador:    datos.Desarrollador,
		FechaLanzamiento: fecha,
		Estado:           "ACTIVO",
	}

	if err := s.Repositorio.CrearVideojuego(
		videojuego,
	); err != nil {

		return nil, err
	}

	return videojuego, nil
}

func (s *CatalogService) ListarCategorias() (
	[]models.Categoria,
	error,
) {

	return s.Repositorio.ListarCategorias()
}

func (s *CatalogService) ObtenerCategoria(
	idCategoria int,
) (*models.Categoria, error) {

	return s.Repositorio.ObtenerCategoria(
		idCategoria,
	)
}

func (s *CatalogService) CrearCategoria(
	datos models.CrearCategoriaInput,
) (*models.Categoria, error) {

	datos.Nombre = strings.TrimSpace(
		datos.Nombre,
	)

	if datos.Nombre == "" {
		return nil, ErrDatosInvalidos
	}

	categoria := &models.Categoria{
		Nombre:      datos.Nombre,
		Descripcion: datos.Descripcion,
	}

	if err := s.Repositorio.CrearCategoria(
		categoria,
	); err != nil {

		return nil, err
	}

	return categoria, nil
}

func (s *CatalogService) ListarPlataformas() (
	[]models.Plataforma,
	error,
) {

	return s.Repositorio.ListarPlataformas()
}

func (s *CatalogService) ObtenerPlataforma(
	idPlataforma int,
) (*models.Plataforma, error) {

	return s.Repositorio.ObtenerPlataforma(
		idPlataforma,
	)
}

func (s *CatalogService) CrearPlataforma(
	datos models.CrearPlataformaInput,
) (*models.Plataforma, error) {

	datos.Nombre = strings.TrimSpace(
		datos.Nombre,
	)

	if datos.Nombre == "" {
		return nil, ErrDatosInvalidos
	}

	plataforma := &models.Plataforma{
		Nombre:      datos.Nombre,
		Descripcion: datos.Descripcion,
	}

	if err := s.Repositorio.CrearPlataforma(
		plataforma,
	); err != nil {

		return nil, err
	}

	return plataforma, nil
}

func (s *CatalogService) AsignarPlataforma(
	datos models.CrearVideojuegoPlataformaInput,
) (*models.VideojuegoPlataforma, error) {

	if datos.IDVideojuego == "" ||
		datos.IDPlataforma <= 0 {

		return nil, ErrDatosInvalidos
	}

	if datos.Precio < 0 {
		return nil, ErrPrecioInvalido
	}

	if datos.Stock < 0 {
		return nil, ErrStockInvalido
	}

	_, err := s.Repositorio.ObtenerVideojuego(
		datos.IDVideojuego,
	)

	if err != nil {
		return nil, err
	}

	_, err = s.Repositorio.ObtenerPlataforma(
		datos.IDPlataforma,
	)

	if err != nil {
		return nil, err
	}

	presentacion :=
		&models.VideojuegoPlataforma{
			IDVideojuego: datos.IDVideojuego,
			IDPlataforma: datos.IDPlataforma,
			Precio:       datos.Precio,
			Stock:        datos.Stock,
			Activo:       true,
		}

	if err :=
		s.Repositorio.CrearVideojuegoPlataforma(
			presentacion,
		); err != nil {

		return nil, err
	}

	plataforma, _ :=
		s.Repositorio.ObtenerPlataforma(
			datos.IDPlataforma,
		)

	if plataforma != nil {
		presentacion.Plataforma = *plataforma
	}

	return presentacion, nil
}

func (s *CatalogService) ObtenerVideojuegoPlataforma(
	id string,
) (*models.VideojuegoPlataforma, error) {

	return s.Repositorio.
		ObtenerVideojuegoPlataforma(id)
}

func (s *CatalogService) ListarPlataformasVideojuego(
	idVideojuego string,
) ([]models.VideojuegoPlataforma, error) {

	return s.Repositorio.
		ListarPlataformasVideojuego(
			idVideojuego,
		)
}

func (s *CatalogService) ConsultarStock(
	id string,
) (int, error) {

	presentacion, err :=
		s.Repositorio.
			ObtenerVideojuegoPlataforma(id)

	if err != nil {
		return 0, err
	}

	return presentacion.Stock, nil
}

func (s *CatalogService) ActualizarStock(
	id string,
	stock int,
) (*models.VideojuegoPlataforma, error) {

	if stock < 0 {
		return nil, ErrStockInvalido
	}

	return s.Repositorio.ActualizarStock(
		id,
		stock,
	)
}
