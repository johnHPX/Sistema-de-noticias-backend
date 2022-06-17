package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/controller/service"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/conteudo"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticia"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

// ==============================
// =========== STORE ============
// ==============================

type storeNoticiaRequest struct {
	Titulo    string              `json:"titulo"`
	Conteudo  []noticia.Conteudos `json:"conteudos"`
	Categoria string              `json:"categoria"`
	MID       string              `json:"mid"`
}

type storeNoticiaResponse struct {
	MID string `json:"mid"`
}

func decodeStoreNoticiaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(storeNoticiaRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func makeStoreNoticiaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*storeNoticiaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}
		service := service.NewNoticiaService()
		var c []conteudo.Conteudo
		for _, v := range req.Conteudo {
			c = append(c, conteudo.Conteudo{
				Subtitulo: v.SubTitulo,
				Texto:     v.Texto,
			})
		}
		err := service.Store(c, req.Titulo, req.Categoria)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1001, err, req.MID)
		}
		//return data
		return &storeNoticiaResponse{
			MID: req.MID,
		}, nil
	}
}

func StoreNoticiaHandler() http.Handler {
	return httptransport.NewServer(
		makeStoreNoticiaEndPoint(),
		decodeStoreNoticiaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== LIST =============
// ==============================

type listNoticiaRequest struct {
	MID string `json:"-"`
}

type listNoticiaResponse struct {
	Count    int                     `json:"count"`
	Entities []noticia.NoticiaEntity `json:"noticias"`
	MID      string                  `json:"mid"`
}

func decodeListNoticiaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	mid := r.URL.Query().Get("mid")
	dto := &listNoticiaRequest{
		MID: mid,
	}
	return dto, nil
}

func makeListNoticiaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*listNoticiaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1002, errors.New("invalid request"), "na")
		}
		service := service.NewNoticiaService()
		entities, err := service.List()
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1003, err, req.MID)
		}
		//return data
		return &listNoticiaResponse{
			Count:    len(entities),
			Entities: entities,
			MID:      req.MID,
		}, nil
	}
}

func ListNoticiaHandler() http.Handler {
	return httptransport.NewServer(
		makeListNoticiaEndPoint(),
		decodeListNoticiaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}
