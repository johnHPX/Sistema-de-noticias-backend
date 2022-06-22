package service

import (
	"github.com/google/uuid"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/categoria"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/repository"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util/validator"
)

type serviceCategoria interface {
	Store(kind string) error
	List() ([]*categoria.Entity, error)
	Find(cid string) (*categoria.Entity, error)
	Update(id, kind string) error
	Remove(id string) error
}

type categoriaServiceImpl struct{}

func (s *categoriaServiceImpl) Store(kind string) error {
	// validator
	validator := validator.NewValidator()
	if err := validator.CheckIsEmpty(kind, "categoria"); err != nil {
		return err
	}
	if str, err := validator.FormatedInput(kind); err != nil {
		kind = str
		return err
	}
	if err := validator.CheckLen(255, kind); err != nil {
		return err
	}

	cid := uuid.New()
	e := &categoria.Entity{
		CID:  cid.String(),
		Kind: kind,
	}
	rep := repository.NewCategoriaRepository()
	err := rep.Store(e)
	if err != nil {
		return err
	}
	return nil
}

func (s *categoriaServiceImpl) List() ([]*categoria.Entity, error) {
	rep := repository.NewCategoriaRepository()
	entities, err := rep.List()
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (s *categoriaServiceImpl) Find(cid string) (*categoria.Entity, error) {
	rep := repository.NewCategoriaRepository()
	entity, err := rep.Find(cid)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *categoriaServiceImpl) Update(id, kind string) error {
	// validator
	validator := validator.NewValidator()
	if err := validator.CheckIsEmpty(kind, "categoria"); err != nil {
		return err
	}
	str, err := validator.FormatedInput(kind)
	if err != nil {
		return err
	}
	kind = str
	if err := validator.CheckLen(255, kind); err != nil {
		return err
	}

	e := &categoria.Entity{
		CID:  id,
		Kind: kind,
	}
	rep := repository.NewCategoriaRepository()
	err = rep.Update(e)
	if err != nil {
		return err
	}
	return nil
}

func (s *categoriaServiceImpl) Remove(id string) error {
	rep := repository.NewCategoriaRepository()
	err := rep.Remove(id)
	if err != nil {
		return err
	}

	return nil
}

func NewCategoriaService() serviceCategoria {
	return &categoriaServiceImpl{}
}
