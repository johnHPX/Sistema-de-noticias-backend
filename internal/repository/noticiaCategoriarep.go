package repository

import (
	"database/sql"
	"errors"

	noticiacategoria "github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticiaCategoria"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

type noticiaCategoriaRepository interface {
	Store(e *noticiacategoria.Entity) error
	List() ([]*noticiacategoria.Entity, error)
}

type noticiaCategoriaRepositoryImpl struct{}

func (r *noticiaCategoriaRepositoryImpl) scanIterator(rows *sql.Rows) (*noticiacategoria.Entity, error) {
	nid := sql.NullString{}
	cid := sql.NullString{}

	err := rows.Scan(
		&nid,
		&cid,
	)

	if err != nil {
		return nil, err
	}

	entity := new(noticiacategoria.Entity)

	if cid.Valid {
		entity.CID = cid.String
	}
	if nid.Valid {
		entity.NID = nid.String
	}

	return entity, nil
}

func (r *noticiaCategoriaRepositoryImpl) Store(e *noticiacategoria.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := "INSERT INTO tb_noticia_categoria(nid,cid) VALUES ($1,$2)"
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.NID, e.CID)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected != 1 {
		return errors.New("error when registering")
	}

	return nil
}

func (r *noticiaCategoriaRepositoryImpl) List() ([]*noticiacategoria.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT nid,cid FROM tb_noticia_categoria
	`
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*noticiacategoria.Entity
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil

}

func NewNoticiaCategoriaRepository() noticiaCategoriaRepository {
	return &noticiaCategoriaRepositoryImpl{}
}
