package repository

import (
	"database/sql"
	"errors"

	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/conteudo"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

type conteudoRepository interface {
	Store(e *conteudo.Entity) error
	ListByNoticia(nid string) ([]*conteudo.Entity, error)
	Update(e *conteudo.Entity) error
	Remove(id string) error
}

type conteudoRepositoryImpl struct{}

func (r *conteudoRepositoryImpl) scanIterator(rows *sql.Rows) (*conteudo.Entity, error) {
	cid := sql.NullString{}
	subtitulo := sql.NullString{}
	texto := sql.NullString{}
	nid := sql.NullString{}

	err := rows.Scan(
		&cid,
		&subtitulo,
		&texto,
		&nid,
	)

	if err != nil {
		return nil, err
	}

	entity := new(conteudo.Entity)

	if cid.Valid {
		entity.CID = cid.String
	}
	if subtitulo.Valid {
		entity.Subtitulo = subtitulo.String
	}
	if texto.Valid {
		entity.Texto = texto.String
	}
	if nid.Valid {
		entity.NID = nid.String
	}

	return entity, nil
}

func (r *conteudoRepositoryImpl) Store(e *conteudo.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := "INSERT INTO tb_conteudo(id, subtitulo, texto, noticia_nid) VALUES ($1,$2,$3,$4)"
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.CID, e.Subtitulo, e.Texto, e.NID)
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

func (r *conteudoRepositoryImpl) ListByNoticia(nid string) ([]*conteudo.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id,subtitulo,texto,noticia_nid FROM tb_conteudo
		WHERE deleted_at is null and noticia_nid = $1
	`
	rows, err := db.Query(sqlText, nid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*conteudo.Entity
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func (r *conteudoRepositoryImpl) Update(e *conteudo.Entity) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_conteudo set
		subtitulo = $2,
		texto = $3,
		updated_at = now()
	where deleted_at is null and id = $1
	`
	statement, err := db.Prepare(sqlText)
	if err != nil {
		return err
	}
	result, err := statement.Exec(e.CID, e.Subtitulo, e.Texto)
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

func (r *conteudoRepositoryImpl) Remove(id string) error {
	db, err := util.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlText := `
	update tb_conteudo set
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

func NewConteudoRepository() conteudoRepository {
	return &conteudoRepositoryImpl{}
}
