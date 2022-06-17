package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/conteudo"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticia"
	noticiacategoria "github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticiaCategoria"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/repository"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util/validator"
)

type serviceNoticia interface {
	Store(conteudos []conteudo.Conteudo, titulo, categoria string) error
	List() ([]noticia.NoticiaEntity, error)
}

type noticiaServiceimpl struct{}

func (s *noticiaServiceimpl) Store(conteudos []conteudo.Conteudo, titulo, cate string) error {
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
	ok := false
	var categoriaID string
	categoriaRep := repository.NewCategoriaRepository()
	categoriaEntities, err := categoriaRep.List()
	if err != nil {
		return err
	}
	for _, v := range *categoriaEntities {
		if cate == v.Kind {
			ok = true
			categoriaID = v.CID
		}
	}
	if !ok {
		return errors.New("Essa categoria é invalida!")
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
		c := conteudo.Conteudo{
			Subtitulo: v.Subtitulo,
			Texto:     v.Texto,
		}
		conteudoEntity := &conteudo.Entity{
			CID:     cid.String(),
			Contudo: c,
			NID:     nid.String(),
		}
		err := conteudoRep.Store(conteudoEntity)
		if err != nil {
			return err
		}
	}
	notCatRep := repository.NewNoticiaCategoriaRepository()
	notCatEntity := &noticiacategoria.Entity{
		NID: nid.String(),
		CID: categoriaID,
	}
	err = notCatRep.Store(notCatEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *noticiaServiceimpl) List() ([]noticia.NoticiaEntity, error) {
	noticaRep := repository.NewNoticiaRepository()
	n, err := noticaRep.List()
	if err != nil {
		return nil, err
	}
	conteudoRep := repository.NewConteudoRepository()
	con, err := conteudoRep.List()
	if err != nil {
		return nil, err
	}
	notConRep := repository.NewNoticiaCategoriaRepository()
	notCon, err := notConRep.List()
	if err != nil {
		return nil, err
	}
	categoriaRep := repository.NewCategoriaRepository()
	entities := make([]noticia.NoticiaEntity, 0)
	for _, not := range *n {
		contuedosOfNoticia := make([]conteudo.Entity, 0)
		for _, conteudo := range *con {
			if not.NID == conteudo.NID {
				contuedosOfNoticia = append(contuedosOfNoticia, conteudo)
			}
		}
		var categoriaOfNoticia string
		for _, Notcategoria := range *notCon {
			if not.NID == Notcategoria.NID {
				e, err := categoriaRep.Find(Notcategoria.CID)
				if err != nil {
					return nil, err
				}
				categoriaOfNoticia = e.Kind
			}
		}
		c := make([]noticia.Conteudos, 0)
		for _, v := range contuedosOfNoticia {
			c = append(c, noticia.Conteudos{
				SubTitulo: v.Contudo.Subtitulo,
				Texto:     v.Contudo.Texto,
			})
		}
		entities = append(entities, noticia.NoticiaEntity{
			ID:        not.NID,
			Titulo:    not.Titulo,
			Conteudo:  c,
			Categoria: categoriaOfNoticia,
		})
	}

	return entities, nil
}

func NewNoticiaService() serviceNoticia {
	return &noticiaServiceimpl{}
}
