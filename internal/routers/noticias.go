package routers

import (
	"net/http"

	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/controller/resource"
)

var RouterNotice = []Router{
	{
		TokenIsReq: false,
		Path:       "/noticia",
		EndPointer: resource.StoreNoticiaHandler().ServeHTTP,
		Method:     http.MethodPost,
	},
	{
		TokenIsReq: false,
		Path:       "/noticias",
		EndPointer: resource.ListNoticiaHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{titcat}/list",
		EndPointer: resource.ListByTitOrCatNoticiaHandler().ServeHTTP,
		Method:     http.MethodGet,
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{id}/update",
		EndPointer: resource.UpdateNoticiaHandler().ServeHTTP,
		Method:     http.MethodPut,
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{id}/remove",
		EndPointer: resource.RemoveNoticiaHandler().ServeHTTP,
		Method:     http.MethodDelete,
	},
}
