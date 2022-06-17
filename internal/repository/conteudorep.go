package repository

import (
	"database/sql"
	"errors"

	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/model/conteudo"
	"github.com/jhonatasfreitas17/sistemaDeNoticias/internal/util"
)

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
		entity.Contudo.Subtitulo = subtitulo.String
	}
	if texto.Valid {
		entity.Contudo.Texto = texto.String
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
	result, err := statement.Exec(e.CID, e.Contudo.Subtitulo, e.Contudo.Texto, e.NID)
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

func (r *conteudoRepositoryImpl) List() (*[]conteudo.Entity, error) {
	db, err := util.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlText := `
		SELECT id,subtitulo,texto,noticia_nid FROM tb_conteudo
	`
	rows, err := db.Query(sqlText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []conteudo.Entity
	for rows.Next() {
		e, err := r.scanIterator(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *e)
	}

	return &entities, nil

}

func NewConteudoRepository() conteudo.Repository {
	return &conteudoRepositoryImpl{}
}
