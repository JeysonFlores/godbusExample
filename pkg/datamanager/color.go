package datamanager

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Color struct {
	con *sql.DB
}

func (p *Color) CreateTable() error {
	sqlStmt := `
			create table if not exists color 
			(id integer not null primary key, 
			value text,
			palette_id integer);
	`
	_, err := p.con.Exec(sqlStmt)
	return err
}

func (p *Color) Insert(name string, palId string) error {
	tx, err := p.con.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into color(id, value, palette_id) values(NULL, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(name, palId)
	tx.Commit()
	stmt.Close()

	return nil
}
