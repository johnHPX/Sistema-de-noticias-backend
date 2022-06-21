package service

import (
	"github.com/google/uuid"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/conteudo"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticia"
	noticiacategoria "github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticiaCategoria"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/repository"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util/validator"
)

type serviceNoticia interface {
	Store(conteudos []conteudo.Entity, titulo, categoria string) error
	List() ([]*noticia.NoticiaEntity, error)
	ListByTitOrCat(titOrCat string) ([]*noticia.NoticiaEntity, error)
}

type noticiaServiceimpl struct{}

func (s *noticiaServiceimpl) Store(conteudos []conteudo.Entity, titulo, cate string) error {
	// validator
	validator := validator.NewValidator()
	for _, v := range conteudos {
		if err := validator.CheckIsEmpty(v.Subtitulo, "subtitulo"); err != nil {
			return err
		}
		if str, err := validator.FormatedInput(v.Subtitulo); err != nil {
			v.Subtitulo = str
			return err
		}
		if err := validator.CheckLen(255, v.Subtitulo); err != nil {
			return err
		}
		if err := validator.CheckIsEmpty(v.Texto, "texto"); err != nil {
			return err
		}
		if str, err := validator.FormatedInput(v.Texto); err != nil {
			v.Texto = str
			return err
		}
		if err := validator.CheckLen(2024, v.Texto); err != nil {
			return err
		}
	}
	if err := validator.CheckIsEmpty(titulo, "titulo"); err != nil {
		return err
	}
	if str, err := validator.FormatedInput(titulo); err != nil {
		titulo = str
		return err
	}
	if err := validator.CheckLen(255, titulo); err != nil {
		return err
	}
	if err := validator.CheckIsEmpty(cate, "categoria"); err != nil {
		return err
	}
	if str, err := validator.FormatedInput(cate); err != nil {
		titulo = str
		return err
	}
	if err := validator.CheckLen(255, cate); err != nil {
		return err
	}

	// verific if categoria exists
	categoriaRep := repository.NewCategoriaRepository()
	categoriaEntity, err := categoriaRep.FindByKind(cate)
	if err != nil {
		return err
	}

	noticaRep := repository.NewNoticiaRepository()
	nid := uuid.New()
	noticiaEntity := &noticia.Entity{
		NID:    nid.String(),
		Titulo: titulo,
	}
	err = noticaRep.Store(noticiaEntity)
	if err != nil {
		return err
	}
	conteudoRep := repository.NewConteudoRepository()
	for _, v := range conteudos {
		cid := uuid.New()
		conteudoEntity := &conteudo.Entity{
			CID:       cid.String(),
			Subtitulo: v.Subtitulo,
			Texto:     v.Texto,
			NID:       nid.String(),
		}
		err := conteudoRep.Store(conteudoEntity)
		if err != nil {
			return err
		}
	}
	notCatRep := repository.NewNoticiaCategoriaRepository()
	notCatEntity := &noticiacategoria.Entity{
		NID: nid.String(),
		CID: categoriaEntity.CID,
	}
	err = notCatRep.Store(notCatEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *noticiaServiceimpl) List() ([]*noticia.NoticiaEntity, error) {
	noticaRep := repository.NewNoticiaRepository()
	entities, err := noticaRep.List()
	if err != nil {
		return nil, err
	}
	conteudoRep := repository.NewConteudoRepository()
	for indice, value := range entities {
		conteudoEntities, err := conteudoRep.ListByNoticia(value.ID)
		if err != nil {
			return nil, err
		}
		c := make([]noticia.Conteudos, 0)
		for _, v := range conteudoEntities {
			c = append(c, noticia.Conteudos{
				SubTitulo: v.Subtitulo,
				Texto:     v.Texto,
			})
		}
		entities[indice].Conteudo = c
	}

	return entities, nil
}

func (s *noticiaServiceimpl) ListByTitOrCat(titOrCat string) ([]*noticia.NoticiaEntity, error) {
	noticaRep := repository.NewNoticiaRepository()
	entities, err := noticaRep.ListByTitOrCat(titOrCat)
	if err != nil {
		return nil, err
	}
	conteudoRep := repository.NewConteudoRepository()
	for indice, value := range entities {
		conteudoEntities, err := conteudoRep.ListByNoticia(value.ID)
		if err != nil {
			return nil, err
		}
		c := make([]noticia.Conteudos, 0)
		for _, v := range conteudoEntities {
			c = append(c, noticia.Conteudos{
				SubTitulo: v.Subtitulo,
				Texto:     v.Texto,
			})
		}
		entities[indice].Conteudo = c
	}

	return entities, nil
}

func NewNoticiaService() serviceNoticia {
	return &noticiaServiceimpl{}
}
