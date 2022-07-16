package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/johnHPX/sistemaDeNoticias/internal/model/conteudo"
	"github.com/johnHPX/sistemaDeNoticias/internal/model/noticia"
	noticiacategoria "github.com/johnHPX/sistemaDeNoticias/internal/model/noticiaCategoria"
	"github.com/johnHPX/sistemaDeNoticias/internal/repository"
	"github.com/johnHPX/sistemaDeNoticias/internal/util/validator"
)

type serviceNoticia interface {
	Store(conteudos []conteudo.Entity, titulo, categoria string) error
	List() ([]*noticia.NoticiaEntity, error)
	ListByTitOrCat(titOrCat string) ([]*noticia.NoticiaEntity, error)
	FindById(id string) (*noticia.NoticiaEntity, error)
	Update(conteudos []conteudo.Entity, NID, titulo, categoria string) error
	Remove(id string) error
}

type noticiaServiceimpl struct{}

func (s *noticiaServiceimpl) Store(conteudos []conteudo.Entity, titulo, cate string) error {
	// validator
	validator := validator.NewValidator()
	conteudosFormated := make([]conteudo.Entity, 0)
	for _, v := range conteudos {
		if err := validator.CheckIsEmpty(v.Subtitulo, "subtitulo"); err != nil {
			return err
		}
		str, err := validator.FormatedInput(v.Subtitulo)
		if err != nil {
			return err
		}
		v.Subtitulo = str
		if err := validator.CheckLen(255, v.Subtitulo); err != nil {
			return err
		}
		if err := validator.CheckIsEmpty(v.Texto, "texto"); err != nil {
			return err
		}
		str, err = validator.FormatedInput(v.Texto)
		if err != nil {
			return err
		}
		v.Texto = str
		if err := validator.CheckLen(5000, v.Texto); err != nil {
			return err
		}
		conteudosFormated = append(conteudosFormated, conteudo.Entity{
			CID:       v.CID,
			Subtitulo: v.Subtitulo,
			Texto:     v.Texto,
			NID:       v.NID,
		})
	}
	if err := validator.CheckIsEmpty(titulo, "titulo"); err != nil {
		return err
	}
	str, err := validator.FormatedInput(titulo)
	if err != nil {
		titulo = str
		return err
	}
	if err := validator.CheckLen(255, titulo); err != nil {
		return err
	}
	if err := validator.CheckIsEmpty(cate, "categoria"); err != nil {
		return err
	}
	str, err = validator.FormatedInput(cate)
	if err != nil {
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
	for _, v := range conteudosFormated {
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
	notCat := uuid.New()
	notCatRep := repository.NewNoticiaCategoriaRepository()
	notCatEntity := &noticiacategoria.Entity{
		ID:  notCat.String(),
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
				CID:       v.CID,
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
				CID:       v.CID,
				SubTitulo: v.Subtitulo,
				Texto:     v.Texto,
			})
		}
		entities[indice].Conteudo = c
	}

	return entities, nil
}

func (s *noticiaServiceimpl) FindById(id string) (*noticia.NoticiaEntity, error) {
	noticaRep := repository.NewNoticiaRepository()
	entity, err := noticaRep.FindById(id)
	if err != nil {
		return nil, err
	}
	conteudoRep := repository.NewConteudoRepository()
	conteudoEntities, err := conteudoRep.ListByNoticia(entity.ID)
	if err != nil {
		return nil, err
	}
	c := make([]noticia.Conteudos, 0)
	for _, v := range conteudoEntities {
		c = append(c, noticia.Conteudos{
			CID:       v.CID,
			SubTitulo: v.Subtitulo,
			Texto:     v.Texto,
		})
	}
	entity.Conteudo = c

	return entity, nil
}

func (s *noticiaServiceimpl) Update(conteudos []conteudo.Entity, NID, titulo, cate string) error {
	// validator
	validator := validator.NewValidator()
	conteudosFormated := make([]conteudo.Entity, 0)
	for _, v := range conteudos {
		if err := validator.CheckIsEmpty(v.Subtitulo, "subtitulo"); err != nil {
			return err
		}
		str, err := validator.FormatedInput(v.Subtitulo)
		if err != nil {
			return err
		}
		v.Subtitulo = str
		if err := validator.CheckLen(255, v.Subtitulo); err != nil {
			return err
		}
		if err := validator.CheckIsEmpty(v.Texto, "texto"); err != nil {
			return err
		}
		str, err = validator.FormatedInput(v.Texto)
		if err != nil {
			return err
		}
		v.Texto = str
		if err := validator.CheckLen(2024, v.Texto); err != nil {
			return err
		}
		conteudosFormated = append(conteudosFormated, conteudo.Entity{
			CID:       v.CID,
			Subtitulo: v.Subtitulo,
			Texto:     v.Texto,
			NID:       v.NID,
		})
	}
	if err := validator.CheckIsEmpty(titulo, "titulo"); err != nil {
		return err
	}
	str, err := validator.FormatedInput(titulo)
	if err != nil {
		titulo = str
		return err
	}
	if err := validator.CheckLen(255, titulo); err != nil {
		return err
	}
	if err := validator.CheckIsEmpty(cate, "categoria"); err != nil {
		return err
	}
	str, err = validator.FormatedInput(cate)
	if err != nil {
		titulo = str
		return err
	}
	if err := validator.CheckLen(255, cate); err != nil {
		return err
	}

	// verific if categoria exists
	categoriaRep := repository.NewCategoriaRepository()
	_, err = categoriaRep.FindByKind(cate)
	if err != nil {
		return err
	}

	// conteudo repository
	conteudoRep := repository.NewConteudoRepository()
	// verific if id of contedos exits
	conteudosEntities, err := conteudoRep.ListByNoticia(NID)
	if err != nil {
		return err
	}
	var count int
	for _, v := range conteudosEntities {
		for _, i := range conteudosFormated {
			if v.CID == i.CID {
				count++
			}
		}
	}
	if count != len(conteudosFormated) {
		return errors.New("Você passou um Conteudo que não está presente nessa noticia.")
	}

	// update conteudos
	for _, v := range conteudosFormated {
		cEntity := conteudo.Entity{
			CID:       v.CID,
			Subtitulo: v.Subtitulo,
			Texto:     v.Texto,
			NID:       NID,
		}
		err := conteudoRep.Update(&cEntity)
		if err != nil {
			return err
		}
	}

	// noticia repository
	noticiaRep := repository.NewNoticiaRepository()
	eEntity := noticia.Entity{
		NID:    NID,
		Titulo: titulo,
	}
	err = noticiaRep.Update(&eEntity)
	if err != nil {
		return err
	}

	// noticiaCategoria repository
	noticiaCategoriaRep := repository.NewNoticiaCategoriaRepository()
	NCEntity, err := noticiaCategoriaRep.FindByNID(NID)
	if err != nil {
		return err
	}

	// categoria repository
	nEntity, err := categoriaRep.FindByKind(cate)
	if err != nil {
		return err
	}

	NCEntity.CID = nEntity.CID
	// update
	err = noticiaCategoriaRep.Update(NCEntity)
	if err != nil {
		return err
	}
	return nil
}

func (s *noticiaServiceimpl) Remove(id string) error {
	conteudoRep := repository.NewConteudoRepository()
	conteudosEntities, err := conteudoRep.ListByNoticia(id)
	if err != nil {
		return err
	}
	for _, v := range conteudosEntities {
		err := conteudoRep.Remove(v.CID)
		if err != nil {
			return err
		}
	}

	noticiaCategoriaRep := repository.NewNoticiaCategoriaRepository()
	NCEntity, err := noticiaCategoriaRep.FindByNID(id)
	if err != nil {
		return err
	}
	err = noticiaCategoriaRep.Remove(NCEntity.ID)
	if err != nil {
		return err
	}

	noticiaRep := repository.NewNoticiaRepository()
	err = noticiaRep.Remove(id)
	if err != nil {
		return err
	}

	return nil
}

func NewNoticiaService() serviceNoticia {
	return &noticiaServiceimpl{}
}
