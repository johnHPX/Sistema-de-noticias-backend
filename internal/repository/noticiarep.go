package repository

import (
	"database/sql"
	"errors"

	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticia"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

type noticiaRepositoryImpl struct{}

func (r *noticiaRepositoryImpl) scanIterator(rows *sql.Rows) (*noticia.Entity, error) {
	nid := sql.NullString{}
	titulo := sql.NullString{}

	err := rows.Scan(
		&nid,
		&titulo,
	)

	if err != nil {
		return nil, err
	}

	entity := new(noticia.Entity)

	if nid.Valid {
		entity.NID = nid.String
	}
	if titulo.Valid {
		entity.Titulo = titulo.String
	}

	return entity, nil
}

func (r *noticiaRepositoryImpl) Store(e *noticia.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := "INSERT INTO tb_noticia(id,titulo) VALUES ($1,$2)"
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.NID, e.Titulo)
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

func (r *noticiaRepositoryImpl) List() (*[]noticia.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id, titulo FROM tb_noticia
	`
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []noticia.Entity
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *e)
	}

	return &entities, nil

}

func NewNoticiaRepository() noticia.Repository {
	return &noticiaRepositoryImpl{}
}
