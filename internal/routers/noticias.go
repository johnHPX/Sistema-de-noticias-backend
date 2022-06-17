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
		Path:       "/noticia/{titCat}",
		EndPointer: nil,
		Method:     "GET",
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{nid}",
		EndPointer: nil,
		Method:     "PUT",
	},
	{
		TokenIsReq: false,
		Path:       "/noticia/{nid}",
		EndPointer: nil,
		Method:     "DELETE",
	},
}
