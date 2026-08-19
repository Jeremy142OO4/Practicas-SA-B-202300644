package graphqlapi

import (
	"github.com/graphql-go/graphql"

	"p4/catalog-ms/models"
	"p4/catalog-ms/services"
)

func CrearEsquema(
	servicio *services.CatalogService,
) (graphql.Schema, error) {

	resolver := NuevoResolver(
		servicio,
	)

	tipoCategoria :=
		graphql.NewObject(
			graphql.ObjectConfig{

				Name: "Categoria",

				Fields: graphql.Fields{

					"idCategoria": &graphql.Field{
						Type: graphql.Int,
					},

					"nombre": &graphql.Field{
						Type: graphql.String,
					},

					"descripcion": &graphql.Field{
						Type: graphql.String,
					},
				},
			},
		)

	tipoPlataforma :=
		graphql.NewObject(
			graphql.ObjectConfig{

				Name: "Plataforma",

				Fields: graphql.Fields{

					"idPlataforma": &graphql.Field{
						Type: graphql.Int,
					},

					"nombre": &graphql.Field{
						Type: graphql.String,
					},

					"descripcion": &graphql.Field{
						Type: graphql.String,
					},
				},
			},
		)

	tipoPresentacion :=
		graphql.NewObject(
			graphql.ObjectConfig{

				Name: "VideojuegoPlataforma",

				Fields: graphql.Fields{

					"idVideojuegoPlataforma": &graphql.Field{
						Type: graphql.ID,
					},

					"idVideojuego": &graphql.Field{
						Type: graphql.ID,
					},

					"idPlataforma": &graphql.Field{
						Type: graphql.Int,
					},

					"precio": &graphql.Field{
						Type: graphql.Float,
					},

					"stock": &graphql.Field{
						Type: graphql.Int,
					},

					"activo": &graphql.Field{
						Type: graphql.Boolean,
					},

					"plataforma": &graphql.Field{

						Type: tipoPlataforma,

						Resolve: func(
							p graphql.ResolveParams,
						) (interface{}, error) {

							switch presentacion :=
								p.Source.(type) {

							case models.VideojuegoPlataforma:
								return presentacion.Plataforma, nil

							case *models.VideojuegoPlataforma:
								return presentacion.Plataforma, nil
							}

							return nil, nil
						},
					},
				},
			},
		)

	tipoVideojuego :=
		graphql.NewObject(
			graphql.ObjectConfig{

				Name: "Videojuego",

				Fields: graphql.Fields{

					"idVideojuego": &graphql.Field{
						Type: graphql.ID,
					},

					"titulo": &graphql.Field{
						Type: graphql.String,
					},

					"descripcion": &graphql.Field{
						Type: graphql.String,
					},

					"idCategoria": &graphql.Field{
						Type: graphql.Int,
					},

					"desarrollador": &graphql.Field{
						Type: graphql.String,
					},

					"fechaLanzamiento": &graphql.Field{

						Type: graphql.String,

						Resolve: func(
							p graphql.ResolveParams,
						) (interface{}, error) {

							switch videojuego :=
								p.Source.(type) {

							case models.Videojuego:

								return videojuego.
									FechaLanzamiento.
									Format("2006-01-02"), nil

							case *models.Videojuego:

								return videojuego.
									FechaLanzamiento.
									Format("2006-01-02"), nil
							}

							return "", nil
						},
					},

					"estado": &graphql.Field{
						Type: graphql.String,
					},

					"categoria": &graphql.Field{

						Type: tipoCategoria,

						Resolve: func(
							p graphql.ResolveParams,
						) (interface{}, error) {

							var idCategoria int

							switch videojuego :=
								p.Source.(type) {

							case models.Videojuego:
								idCategoria =
									videojuego.IDCategoria

							case *models.Videojuego:
								idCategoria =
									videojuego.IDCategoria
							}

							return servicio.
								ObtenerCategoria(
									idCategoria,
								)
						},
					},

					"plataformas": &graphql.Field{

						Type: graphql.NewList(
							tipoPresentacion,
						),

						Resolve: func(
							p graphql.ResolveParams,
						) (interface{}, error) {

							var idVideojuego string

							switch videojuego :=
								p.Source.(type) {

							case models.Videojuego:
								idVideojuego =
									videojuego.IDVideojuego

							case *models.Videojuego:
								idVideojuego =
									videojuego.IDVideojuego
							}

							return servicio.
								ListarPlataformasVideojuego(
									idVideojuego,
								)
						},
					},
				},
			},
		)

	query :=
		graphql.NewObject(
			graphql.ObjectConfig{

				Name: "Query",

				Fields: graphql.Fields{

					"videojuegos": &graphql.Field{

						Type: graphql.NewList(
							tipoVideojuego,
						),

						Resolve: resolver.
							ListarVideojuegos,
					},

					"videojuego": &graphql.Field{

						Type: tipoVideojuego,

						Args: graphql.FieldConfigArgument{

							"id": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.ID,
								),
							},
						},

						Resolve: resolver.
							ObtenerVideojuego,
					},

					"categorias": &graphql.Field{

						Type: graphql.NewList(
							tipoCategoria,
						),

						Resolve: resolver.
							ListarCategorias,
					},

					"plataformas": &graphql.Field{

						Type: graphql.NewList(
							tipoPlataforma,
						),

						Resolve: resolver.
							ListarPlataformas,
					},

					"stock": &graphql.Field{

						Type: graphql.Int,

						Args: graphql.FieldConfigArgument{

							"id": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.ID,
								),
							},
						},

						Resolve: resolver.
							ConsultarStock,
					},

					"videojuegoPlataforma": &graphql.Field{

						Type: tipoPresentacion,

						Args: graphql.FieldConfigArgument{

							"id": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.ID,
								),
							},
						},

						Resolve: resolver.
							ObtenerPresentacion,
					},
				},
			},
		)

	mutation :=
		graphql.NewObject(
			graphql.ObjectConfig{

				Name: "Mutation",

				Fields: graphql.Fields{

					"crearCategoria": &graphql.Field{

						Type: tipoCategoria,

						Args: graphql.FieldConfigArgument{

							"nombre": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.String,
								),
							},

							"descripcion": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},

						Resolve: resolver.
							CrearCategoria,
					},

					"crearPlataforma": &graphql.Field{

						Type: tipoPlataforma,

						Args: graphql.FieldConfigArgument{

							"nombre": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.String,
								),
							},

							"descripcion": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},

						Resolve: resolver.
							CrearPlataforma,
					},

					"crearVideojuego": &graphql.Field{

						Type: tipoVideojuego,

						Args: graphql.FieldConfigArgument{

							"titulo": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.String,
								),
							},

							"descripcion": &graphql.ArgumentConfig{
								Type: graphql.String,
							},

							"idCategoria": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.Int,
								),
							},

							"desarrollador": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.String,
								),
							},

							"fechaLanzamiento": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.String,
								),
							},
						},

						Resolve: resolver.
							CrearVideojuego,
					},

					"asignarPlataforma": &graphql.Field{

						Type: tipoPresentacion,

						Args: graphql.FieldConfigArgument{

							"idVideojuego": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.ID,
								),
							},

							"idPlataforma": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.Int,
								),
							},

							"precio": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.Float,
								),
							},

							"stock": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.Int,
								),
							},
						},

						Resolve: resolver.
							AsignarPlataforma,
					},

					"actualizarStock": &graphql.Field{

						Type: tipoPresentacion,

						Args: graphql.FieldConfigArgument{

							"id": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.ID,
								),
							},

							"stock": &graphql.ArgumentConfig{
								Type: graphql.NewNonNull(
									graphql.Int,
								),
							},
						},

						Resolve: resolver.
							ActualizarStock,
					},
				},
			},
		)

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    query,
			Mutation: mutation,
		},
	)
}
