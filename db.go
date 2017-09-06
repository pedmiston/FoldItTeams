package main

import (
	"database/sql"
	"fmt"
	"strings"
)

func createTables() error {
	sqlStmt := `
  create table scores(filename text not null primary key, puzzle_id integer, user_id integer, group_id integer, score float);
  create table actions(filename text, action text, count integer);
  create table history(filename text, ix integer, id text);
  create table ranks(filename text not null primary key, type text, rank integer);
  `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func prepareStatement(tx *sql.Tx, tblName string, fields []string) (*sql.Stmt, error) {
	q := make([]string, len(fields))
	for i := range fields {
		q[i] = "?"
	}

	msg := fmt.Sprintf("insert into %s(%s) values(%s)",
		tblName, strings.Join(fields, ","), strings.Join(q, ","))

	stmt, err := tx.Prepare(msg)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}
