package routers

import (
	"net/http"

	"github.com/johnHPX/sistemaDeNoticias/internal/controller/resource"
)

var RouterCategoria = []Router{
	{
		TokenIsReq: false,
		Path:       "/categoria",
		EndPointer: resource.StoreCategoriaHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/categorias",
		EndPointer: resource.ListCategoriaHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/categoria/{id}/find",
		EndPointer: resource.FindCategoriaHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/categoria/{id}/update",
		EndPointer: resource.UpdateCategoriaHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: false,
		Path:       "/categoria/{id}/remove",
		EndPointer: resource.RemoveCategoriaHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
