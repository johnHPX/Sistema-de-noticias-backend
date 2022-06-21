package routers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

type Router struct {
	TokenIsReq bool
	Path       string
	EndPointer http.HandlerFunc
	Method     string
}

type WebService interface {
	Init()
	GetRouters() http.Handler
}

type webServiceImpl struct {
	Router *mux.Router
	ctx    context.Context
}

func (s *webServiceImpl) Init() {
	s.configuration()
}

func (s *webServiceImpl) GetRouters() http.Handler {
	return s.Router
}

func (s *webServiceImpl) configuration() {
	routers := []Router{}
	routers = append(routers, RouterNotice...)
	routers = append(routers, RouterCategoria...)
	for _, router := range routers {
		if router.TokenIsReq {
			s.Router.HandleFunc(router.Path, util.Logger(util.Authenticate(router.EndPointer))).Methods(router.Method)
		}
		s.Router.HandleFunc(router.Path, util.Logger(router.EndPointer)).Methods(router.Method)
	}
}

func NewWebService(ctx context.Context) WebService {
	return &webServiceImpl{
		Router: mux.NewRouter(),
		ctx:    ctx,
	}
}
