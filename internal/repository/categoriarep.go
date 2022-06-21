package repository

import (
	"database/sql"
	"errors"

	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/categoria"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

type categoriaRepository interface {
	Store(e *categoria.Entity) error
	List() ([]*categoria.Entity, error)
	Find(cid string) (*categoria.Entity, error)
	FindByKind(kind string) (*categoria.Entity, error)
	Update(e *categoria.Entity) error
}

type categoriaRepositoryImpl struct{}

func (r *categoriaRepositoryImpl) scanIterator(rows *sql.Rows) (*categoria.Entity, error) {
	cid := sql.NullString{}
	kind := sql.NullString{}

	err := rows.Scan(
		&cid,
		&kind,
	)

	if err != nil {
		return nil, err
	}

	entity := new(categoria.Entity)

	if cid.Valid {
		entity.CID = cid.String
	}
	if kind.Valid {
		entity.Kind = kind.String
	}

	return entity, nil
}

func (r *categoriaRepositoryImpl) Store(e *categoria.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := "INSERT INTO tb_categoria(id,kind) VALUES ($1,$2)"
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.CID, e.Kind)
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

func (r *categoriaRepositoryImpl) List() ([]*categoria.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id,kind FROM tb_categoria WHERE deleted_at is null
	`
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*categoria.Entity
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func (r *categoriaRepositoryImpl) Find(cid string) (*categoria.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id,kind FROM tb_categoria
		WHERE deleted_at is null and id = $1
	`
	row, err := db.Query(sqlText, cid)
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

func (r *categoriaRepositoryImpl) FindByKind(kind string) (*categoria.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id,kind FROM tb_categoria
		WHERE deleted_at is null and kind = $1
	`
	row, err := db.Query(sqlText, kind)
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

	return nil, errors.New("error finding categoria")
}

func (r *categoriaRepositoryImpl) Update(e *categoria.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_categoria set
		kind = $2,
		updated_at = now()
	where deleted_at is null and id = $1
	`
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.CID, e.Kind)
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

func NewCategoriaRepository() categoriaRepository {
	return &categoriaRepositoryImpl{}
}
