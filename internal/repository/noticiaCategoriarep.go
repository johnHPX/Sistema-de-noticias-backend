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
	FindByNID(nid string) (*noticiacategoria.Entity, error)
	Update(e *noticiacategoria.Entity) error
}

type noticiaCategoriaRepositoryImpl struct{}

func (r *noticiaCategoriaRepositoryImpl) scanIterator(rows *sql.Rows) (*noticiacategoria.Entity, error) {
	id := sql.NullString{}
	nid := sql.NullString{}
	cid := sql.NullString{}

	err := rows.Scan(
		&id,
		&nid,
		&cid,
	)

	if err != nil {
		return nil, err
	}

	entity := new(noticiacategoria.Entity)

	if id.Valid {
		entity.ID = id.String
	}

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
	sqlText := "INSERT INTO tb_noticia_categoria(id, nid,cid) VALUES ($1,$2, $3)"
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.ID, e.NID, e.CID)
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
		SELECT id, nid, cid FROM tb_noticia_categoria WHERE deleted_at is null
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

func (r *noticiaCategoriaRepositoryImpl) FindByNID(nid string) (*noticiacategoria.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id, nid,cid FROM tb_noticia_categoria WHERE deleted_at is null and nid = $1
	`
	rows, err := db.Query(sqlText, nid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		return e, nil
	}

	return nil, errors.New("error finding")
}

func (r *noticiaCategoriaRepositoryImpl) Update(e *noticiacategoria.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_noticia_categoria set
		nid = $2,
		cid = $3,
		updated_at = now()
	where deleted_at is null and id = $1
	`
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.ID, e.NID, e.CID)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when updating")
	}

	return nil
}

func NewNoticiaCategoriaRepository() noticiaCategoriaRepository {
	return &noticiaCategoriaRepositoryImpl{}
}
