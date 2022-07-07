package resource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/controller/service"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/conteudo"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticia"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

// ==============================
// =========== STORE ============
// ==============================

type storeNoticiaRequest struct {
	noticia.NoticiaEntity
	MID string `json:"mid"`
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
		var c []conteudo.Entity
		for _, v := range req.Conteudo {
			c = append(c, conteudo.Entity{
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
	Count    int                      `json:"count"`
	Entities []*noticia.NoticiaEntity `json:"noticias"`
	MID      string                   `json:"mid"`
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

// =====================================================
// =========== LIST BY TITULO OR CATEGORIA =============
// =====================================================

type listByTitOrCatNoticiaRequest struct {
	TITCAT string `json:"-"`
	MID    string `json:"-"`
}

type listByTitOrCatNoticiaResponse struct {
	Count    int                     `json:"count"`
	Entities []noticia.NoticiaEntity `json:"noticias"`
	MID      string                  `json:"mid"`
}

func decodeListByTitOrCatNoticiaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	mid := r.URL.Query().Get("mid")
	dto := &listByTitOrCatNoticiaRequest{
		TITCAT: vars["titcat"],
		MID:    mid,
	}
	return dto, nil
}

func makeListByTitOrCatNoticiaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*listByTitOrCatNoticiaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1004, errors.New("invalid request"), "na")
		}
		service := service.NewNoticiaService()
		entities, err := service.ListByTitOrCat(req.TITCAT)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}
		//return data
		return &listNoticiaResponse{
			Count:    len(entities),
			Entities: entities,
			MID:      req.MID,
		}, nil
	}
}

func ListByTitOrCatNoticiaHandler() http.Handler {
	return httptransport.NewServer(
		makeListByTitOrCatNoticiaEndPoint(),
		decodeListByTitOrCatNoticiaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// ===== FIND BY ID  ===========
// ==============================

type findNoticiaRequest struct {
	ID  string `json:"-"`
	MID string `json:"-"`
}

type findNoticiaResponse struct {
	noticia.NoticiaEntity
	MID string `json:"mid"`
}

func decodeFindNoticiaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	mid := r.URL.Query().Get("mid")
	dto := &findNoticiaRequest{
		ID:  vars["id"],
		MID: mid,
	}
	return dto, nil
}

func makeFindNoticiaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*findNoticiaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1006, errors.New("invalid request"), "na")
		}
		service := service.NewNoticiaService()
		entity, err := service.FindById(req.ID)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1007, err, req.MID)
		}
		//return data
		return &findNoticiaResponse{
			NoticiaEntity: *entity,
			MID:           req.MID,
		}, nil
	}
}

func FindNoticiaHandler() http.Handler {
	return httptransport.NewServer(
		makeFindNoticiaEndPoint(),
		decodeFindNoticiaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== UPDATE ===========
// ==============================

type updateNoticiaRequest struct {
	noticia.NoticiaEntity
	MID string `json:"mid"`
}

type updateNoticiaResponse struct {
	MID string `json:"mid"`
}

func decodeUpdateNoticiaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dto := new(updateNoticiaRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.ID = vars["id"]
	return dto, nil
}

func makeUpdateNoticiaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*updateNoticiaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1006, errors.New("invalid request"), "na")
		}
		service := service.NewNoticiaService()
		var c []conteudo.Entity
		for _, v := range req.Conteudo {
			c = append(c, conteudo.Entity{
				CID:       v.CID,
				Subtitulo: v.SubTitulo,
				Texto:     v.Texto,
			})
		}
		err := service.Update(c, req.ID, req.Titulo, req.Categoria)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1007, err, req.MID)
		}
		//return data
		return &updateNoticiaResponse{
			MID: req.MID,
		}, nil
	}
}

func UpdateNoticiaHandler() http.Handler {
	return httptransport.NewServer(
		makeUpdateNoticiaEndPoint(),
		decodeUpdateNoticiaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== REMOVE ===========
// ==============================

type removeNoticiaRequest struct {
	ID  string `json:"-"`
	MID string `json:"-"`
}

type removeNoticiaResponse struct {
	MID string `json:"mid"`
}

func decodeRemoveNoticiaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dto := new(removeNoticiaRequest)
	dto.ID = vars["id"]
	dto.MID = r.URL.Query().Get("mid")
	return dto, nil
}

func makeRemoveNoticiaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*removeNoticiaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1008, errors.New("invalid request"), "na")
		}
		service := service.NewNoticiaService()
		err := service.Remove(req.ID)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1009, err, req.MID)
		}
		//return data
		return &removeNoticiaResponse{
			MID: req.MID,
		}, nil
	}
}

func RemoveNoticiaHandler() http.Handler {
	return httptransport.NewServer(
		makeRemoveNoticiaEndPoint(),
		decodeRemoveNoticiaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}
