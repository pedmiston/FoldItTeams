/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -paths=paths.txt -db=foldit.sqlite

Args:
	paths: A file containing paths to solution files. If a file is not provided,
		paths are expected via stdin.
	db: Location of a sqlite database containing the extracted data.
*/
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func main() {
	paths := flag.String("pdbs", "",
		"A file containing paths to files to process.")
	dbName := flag.String("db", "",
		"The name of the sqlite database that will contain the extracted data.")
	flag.Parse()

	// Create a solution file path scanner
	var src *os.File
	var err error
	if *paths == "" {
		src = os.Stdin
	} else {
		src, err = os.Open(*paths)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer src.Close()
	scanner := bufio.NewScanner(src)

	// Create an interface to the DB
	if *dbName == "" {
		*dbName = "foldit.db"
	}
	os.Remove(*dbName)
	db, err = sql.Open("sqlite3", *dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tables
	sqlStmt := `
	create table scores(filename text not null primary key, puzzle_id integer, user_id integer, group_id integer, score float);
	create table actions(filename text not null primary key, action text, count integer);
	create table history(filename text not null primary key, ix integer, id text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	ExtractSolutionData(scanner, db)
}

// ExtractSolutionData extracts data from pdb files and inserts it into a db.
func ExtractSolutionData(paths *bufio.Scanner, db *sql.DB) {
	var ch chan *Solution

	// Load solutions in a chan
	ch = loadSolutions(paths)

	for solution := range ch {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		stmtScores := solution.prepareScores(tx)
		stmtActions := solution.prepareActions(tx)
		stmtHistory := solution.prepareHistory(tx)

		stmtScores.Exec(solution.getScores())

		for _, d := range solution.getActions() {
			stmtActions.Exec(d)
		}

		for _, d := range solution.getHistory() {
			stmtHistory.Exec(d)
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		stmtScores.Close()
		stmtActions.Close()
		stmtHistory.Close()
	}
}

func loadSolutions(in *bufio.Scanner) chan *Solution {
	out := make(chan *Solution)
	go func() {
		for in.Scan() {
			out <- NewSolution(in.Text())
		}
		close(out)
	}()
	return out
}
