package routers

import "github.com/jhonatasfreitas17/sistemaDeNoticias/internal/controller/resource"

var RouterCategoria = []Router{
	{
		TokenIsReq: false,
		Path:       "/categoria",
		EndPointer: resource.StoreCategoriaHandler().ServeHTTP,
		Method:     "POST",
	},
	{
		TokenIsReq: false,
		Path:       "/categorias",
		EndPointer: resource.ListCategoriaHandler().ServeHTTP,
		Method:     "GET",
	},
	{
		TokenIsReq: false,
		Path:       "/categoria/{id}",
		EndPointer: resource.FindCategoriaHandler().ServeHTTP,
		Method:     "GET",
	},
	{
		TokenIsReq: false,
		Path:       "/categoria/{id}",
		EndPointer: resource.UpdateCategoriaHandler().ServeHTTP,
		Method:     "PUT",
	},
	{
		TokenIsReq: false,
		Path:       "/categoria/{id}",
		EndPointer: resource.RemoveCategoriaHandler().ServeHTTP,
		Method:     "DELETE",
	},
}
