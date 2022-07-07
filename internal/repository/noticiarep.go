package repository

import (
	"database/sql"
	"errors"

	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/noticia"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

type noticiaRepository interface {
	Store(e *noticia.Entity) error
	List() ([]*noticia.NoticiaEntity, error)
	ListByTitOrCat(titCat string) ([]*noticia.NoticiaEntity, error)
	FindById(id string) (*noticia.NoticiaEntity, error)
	Update(e *noticia.Entity) error
	Remove(id string) error
}

type noticiaRepositoryImpl struct{}

func (r *noticiaRepositoryImpl) scanIterator(rows *sql.Rows) (*noticia.NoticiaEntity, error) {
	nid := sql.NullString{}
	titulo := sql.NullString{}
	categoria := sql.NullString{}

	err := rows.Scan(
		&nid,
		&titulo,
		&categoria,
	)

	if err != nil {
		return nil, err
	}

	entity := new(noticia.NoticiaEntity)

	if nid.Valid {
		entity.ID = nid.String
	}
	if titulo.Valid {
		entity.Titulo = titulo.String
	}
	if categoria.Valid {
		entity.Categoria = categoria.String
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

func (r *noticiaRepositoryImpl) List() ([]*noticia.NoticiaEntity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
	select n.id,
		n.titulo,
		ca.kind
	from tb_noticia  n
	INNER JOIN tb_noticia_categoria nc ON nc.nid = n.id
	INNER JOIN tb_categoria ca ON ca.id = nc.cid
	where n.deleted_at is null and ca.deleted_at is null
	`
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*noticia.NoticiaEntity
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil

}

func (r *noticiaRepositoryImpl) ListByTitOrCat(titCat string) ([]*noticia.NoticiaEntity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
	select n.id,
		n.titulo,
		ca.kind
	from tb_noticia  n
	INNER JOIN tb_noticia_categoria nc ON nc.nid = n.id
	INNER JOIN tb_categoria ca ON ca.id = nc.cid
	where n.deleted_at is null and ca.deleted_at is null
	and (n.titulo like $1 or ca.kind like $2)
	`

	v := "%" + titCat + "%"

	rows, err := db.Query(sqlText, v, v)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*noticia.NoticiaEntity, 0)
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func (r *noticiaRepositoryImpl) FindById(id string) (*noticia.NoticiaEntity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
	select n.id,
		n.titulo,
		ca.kind
	from tb_noticia  n
	INNER JOIN tb_noticia_categoria nc ON nc.nid = n.id
	INNER JOIN tb_categoria ca ON ca.id = nc.cid
	where n.deleted_at is null and ca.deleted_at is null and n.id = $1
	`
	row, err := db.Query(sqlText, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		e, err := r.scanIterator(row)
		if err != nil {
			return nil, err
		}
		return e, nil
	}

	return nil, errors.New("error finding")
}

func (r *noticiaRepositoryImpl) Update(e *noticia.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_noticia set
		titulo = $2,
		updated_at = now()
	where deleted_at is null and id = $1
	`
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
		return errors.New("error when updating")
	}

	return nil
}

func (r *noticiaRepositoryImpl) Remove(id string) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_noticia set
		deleted_at = now()
	where id = $1
	`
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(id)
	if err != nil {
		return err
	}
	defer statement.Close()

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return errors.New("error when deleting")
	}

	return nil
}

func NewNoticiaRepository() noticiaRepository {
	return &noticiaRepositoryImpl{}
}
