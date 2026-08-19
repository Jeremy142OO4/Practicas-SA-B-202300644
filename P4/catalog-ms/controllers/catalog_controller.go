package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type CatalogController struct {
	Esquema graphql.Schema
}

type PeticionGraphQL struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

func NuevoCatalogController(
	esquema graphql.Schema,
) *CatalogController {

	return &CatalogController{
		Esquema: esquema,
	}
}

func (c *CatalogController) EjecutarGraphQL(
	ctx *fiber.Ctx,
) error {

	var peticion PeticionGraphQL

	if err := ctx.BodyParser(
		&peticion,
	); err != nil {

		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": "peticion inválida",
			},
		)
	}

	if peticion.Query == "" {

		return ctx.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": "query requerida",
			},
		)
	}

	resultado := graphql.Do(
		graphql.Params{

			Schema: c.Esquema,

			RequestString: peticion.Query,

			VariableValues: peticion.Variables,

			OperationName: peticion.OperationName,
		},
	)

	return ctx.JSON(resultado)
}
