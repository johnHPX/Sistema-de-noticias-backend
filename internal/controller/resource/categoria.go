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
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/categoria"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

// ==============================
// =========== STORE ============
// ==============================

type storeCategoriaRequest struct {
	Kind string `json:"kind"`
	MID  string `json:"mid"`
}

type storeCategoriaResponse struct {
	MID string `json:"mid"`
}

func decodeStoreCategoriaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	dto := new(storeCategoriaRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func makeStoreCategoriandPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*storeCategoriaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1000, errors.New("invalid request"), "na")
		}
		svc := service.NewCategoriaService()
		err := svc.Store(req.Kind)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1001, err, req.MID)
		}
		//return data
		return &storeCategoriaResponse{
			MID: req.MID,
		}, nil
	}
}

func StoreCategoriaHandler() http.Handler {
	return httptransport.NewServer(
		makeStoreCategoriandPoint(),
		decodeStoreCategoriaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== LIST =============
// ==============================

type listCategoriaRequest struct {
	MID string `json:"-"`
}

type listCategoriaResponse struct {
	Count    int                 `json:"count"`
	Entities []*categoria.Entity `json:"categorias"`
	MID      string              `json:"mid"`
}

func decodeListCategoriaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	mid := r.URL.Query().Get("mid")
	dto := &listCategoriaRequest{
		MID: mid,
	}
	return dto, nil
}

func makeListCategoriaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*listCategoriaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1002, errors.New("invalid request"), "na")
		}
		service := service.NewCategoriaService()
		entities, err := service.List()
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1003, err, req.MID)
		}
		//return data
		return &listCategoriaResponse{
			Count:    len(entities),
			Entities: entities,
			MID:      req.MID,
		}, nil
	}
}

func ListCategoriaHandler() http.Handler {
	return httptransport.NewServer(
		makeListCategoriaEndPoint(),
		decodeListCategoriaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== FIND =============
// ==============================

type findCategoriaRequest struct {
	CID string `json:"-"`
	MID string `json:"-"`
}

type findCategoriaResponse struct {
	CID  string `json:"id"`
	Kind string `json:"kind"`
	MID  string `json:"mid"`
}

func decodeFindCategoriaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	mid := r.URL.Query().Get("mid")
	dto := &findCategoriaRequest{
		CID: vars["id"],
		MID: mid,
	}
	return dto, nil
}

func makeFindCategoriaEndPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*findCategoriaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1004, errors.New("invalid request"), "na")
		}
		svc := service.NewCategoriaService()
		entity, err := svc.Find(req.CID)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1005, err, req.MID)
		}
		//return data
		return &findCategoriaResponse{
			CID:  entity.CID,
			Kind: entity.Kind,
			MID:  req.MID,
		}, nil
	}
}

func FindCategoriaHandler() http.Handler {
	return httptransport.NewServer(
		makeFindCategoriaEndPoint(),
		decodeFindCategoriaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== UPDATE ===========
// ==============================

type updateCategoriaRequest struct {
	ID   string
	Kind string `json:"kind"`
	MID  string `json:"mid"`
}

type updateCategoriaResponse struct {
	MID string `json:"mid"`
}

func decodeUpdateategoriaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dto := new(updateCategoriaRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dto)
	if err != nil {
		return nil, err
	}
	dto.ID = vars["id"]
	return dto, nil
}

func makeUpdateCategoriandPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*updateCategoriaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1006, errors.New("invalid request"), "na")
		}
		svc := service.NewCategoriaService()
		err := svc.Update(req.ID, req.Kind)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1007, err, req.MID)
		}
		//return data
		return &updateCategoriaResponse{
			MID: req.MID,
		}, nil
	}
}

func UpdateCategoriaHandler() http.Handler {
	return httptransport.NewServer(
		makeUpdateCategoriandPoint(),
		decodeUpdateategoriaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}

// ==============================
// =========== REMOVE ===========
// ==============================

type removeCategoriaRequest struct {
	ID  string `json:"-"`
	MID string `json:"-"`
}

type removeCategoriaResponse struct {
	MID string `json:"mid"`
}

func decodeRemoveategoriaRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dto := new(removeCategoriaRequest)
	dto.ID = vars["id"]
	return dto, nil
}

func makeRemoveCategoriandPoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// retrieve request data
		req, ok := request.(*removeCategoriaRequest)
		if !ok {
			return nil, util.CreateHttpErrorResponse(http.StatusBadRequest, 1008, errors.New("invalid request"), "na")
		}
		svc := service.NewCategoriaService()
		err := svc.Remove(req.ID)
		if err != nil {
			return nil, util.CreateHttpErrorResponse(http.StatusInternalServerError, 1009, err, req.MID)
		}
		//return data
		return &updateCategoriaResponse{
			MID: req.MID,
		}, nil
	}
}

func RemoveCategoriaHandler() http.Handler {
	return httptransport.NewServer(
		makeRemoveCategoriandPoint(),
		decodeRemoveategoriaRequest,
		util.EncodeResponse,
		httptransport.ServerErrorEncoder(util.ErrorEncoder()),
	)
}
