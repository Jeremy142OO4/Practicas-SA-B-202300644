package graphqlapi

import (
	"github.com/graphql-go/graphql"

	"p4/catalog-ms/models"
	"p4/catalog-ms/services"
)

type Resolver struct {
	Servicio *services.CatalogService
}

func NuevoResolver(
	servicio *services.CatalogService,
) *Resolver {

	return &Resolver{
		Servicio: servicio,
	}
}

func (r *Resolver) ListarVideojuegos(
	p graphql.ResolveParams,
) (interface{}, error) {

	return r.Servicio.ListarVideojuegos()
}

func (r *Resolver) ObtenerVideojuego(
	p graphql.ResolveParams,
) (interface{}, error) {

	id := p.Args["id"].(string)

	return r.Servicio.ObtenerVideojuego(id)
}

func (r *Resolver) ListarCategorias(
	p graphql.ResolveParams,
) (interface{}, error) {

	return r.Servicio.ListarCategorias()
}

func (r *Resolver) ListarPlataformas(
	p graphql.ResolveParams,
) (interface{}, error) {

	return r.Servicio.ListarPlataformas()
}

func (r *Resolver) ConsultarStock(
	p graphql.ResolveParams,
) (interface{}, error) {

	id := p.Args["id"].(string)

	return r.Servicio.ConsultarStock(id)
}

func (r *Resolver) ObtenerPresentacion(
	p graphql.ResolveParams,
) (interface{}, error) {

	id := p.Args["id"].(string)

	return r.Servicio.
		ObtenerVideojuegoPlataforma(id)
}

func (r *Resolver) CrearCategoria(
	p graphql.ResolveParams,
) (interface{}, error) {

	datos := models.CrearCategoriaInput{
		Nombre: p.Args["nombre"].(string),
	}

	if descripcion, ok :=
		p.Args["descripcion"].(string); ok {

		datos.Descripcion = descripcion
	}

	return r.Servicio.CrearCategoria(datos)
}

func (r *Resolver) CrearPlataforma(
	p graphql.ResolveParams,
) (interface{}, error) {

	datos := models.CrearPlataformaInput{
		Nombre: p.Args["nombre"].(string),
	}

	if descripcion, ok :=
		p.Args["descripcion"].(string); ok {

		datos.Descripcion = descripcion
	}

	return r.Servicio.CrearPlataforma(datos)
}

func (r *Resolver) CrearVideojuego(
	p graphql.ResolveParams,
) (interface{}, error) {

	datos := models.CrearVideojuegoInput{
		Titulo: p.Args["titulo"].(string),

		IDCategoria: p.Args["idCategoria"].(int),

		Desarrollador: p.Args["desarrollador"].(string),

		FechaLanzamiento: p.Args["fechaLanzamiento"].(string),
	}

	if descripcion, ok :=
		p.Args["descripcion"].(string); ok {

		datos.Descripcion = descripcion
	}

	return r.Servicio.CrearVideojuego(datos)
}

func (r *Resolver) AsignarPlataforma(
	p graphql.ResolveParams,
) (interface{}, error) {

	datos :=
		models.CrearVideojuegoPlataformaInput{

			IDVideojuego: p.Args["idVideojuego"].(string),

			IDPlataforma: p.Args["idPlataforma"].(int),

			Precio: p.Args["precio"].(float64),

			Stock: p.Args["stock"].(int),
		}

	return r.Servicio.AsignarPlataforma(
		datos,
	)
}

func (r *Resolver) ActualizarStock(
	p graphql.ResolveParams,
) (interface{}, error) {

	id := p.Args["id"].(string)

	stock := p.Args["stock"].(int)

	return r.Servicio.ActualizarStock(
		id,
		stock,
	)
}
