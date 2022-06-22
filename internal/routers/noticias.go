package routers

import "github.com/jhonatasfreitas17/sistemaDeNoticias/internal/controller/resource"

var RouterNotice = []Router{
	{
		TokenIsReq: false,
		Path:       "/noticia",
		EndPointer: resource.StoreNoticiaHandler().ServeHTTP,
		Method:     "POST",
	},
	{
		TokenIsReq: false,
		Path:       "/noticias",
		EndPointer: resource.ListNoticiaHandler().ServeHTTP,
		Method:     "GET",
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{titcat}",
		EndPointer: resource.ListByTitOrCatNoticiaHandler().ServeHTTP,
		Method:     "GET",
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{id}",
		EndPointer: resource.UpdateNoticiaHandler().ServeHTTP,
		Method:     "PUT",
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{id}",
		EndPointer: resource.RemoveNoticiaHandler().ServeHTTP,
		Method:     "DELETE",
	},
}
