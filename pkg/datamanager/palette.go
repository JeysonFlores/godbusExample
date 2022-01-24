package datamanager

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Palette struct {
	Con *sql.DB
}

func (p *Palette) CreateTable() error {
	sqlStmt := `
			create table if not exists palette 
			(id integer not null primary key, 
			name text,
			description text);
	`
	_, err := p.Con.Exec(sqlStmt)
	return err
}

func (p *Palette) Insert(name string, desc string) error {
	tx, err := p.Con.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into palette(id, name, description) values(NULL, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(name, desc)
	tx.Commit()
	stmt.Close()

	return nil
}

func (p *Palette) GetById(id int) (string, error) {
	stmt, err := p.Con.Prepare("SELECT name FROM palette WHERE id = ?")
	fmt.Println("query prepared")
	if err != nil {
		fmt.Println(err)
		return "there was an error preparing the query", err
	}

	defer stmt.Close()

	var name string
	err = stmt.QueryRow(id).Scan(&name)
	if err != nil {
		return "there's no palette that matches the given id", err
	}

	return name, nil
}
