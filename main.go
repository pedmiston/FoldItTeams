/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -paths=solution-filepaths.txt -dst=foldit.sqlite

Args:
	paths: A file containing paths to solution files. If a file is not provided,
		paths are expected via stdin.
	dst: Location of a sqlite database containing the extracted data.
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
	paths := flag.String("paths", "",
		"A file containing paths to files to process.")
	dst := flag.String("dst", "", "Location of a sqlite database")
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
	if *dst == "" {
		*dst = "foldit.db"
	}
	os.Remove(*dst)
	db, err = sql.Open("sqlite3", *dst)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tables
	sqlStmt := `
	create table scores(filename text not null primary key, puzzle_id integer, user_id integer, group_id integer, score float);
	create table actions(filename text not null primary key, action text, count integer);
	create table history(filename text not null primary key, ix integer, id text);
	create table ranks(filename text not null primary key, type text, rank integer);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	ExtractSolutionData(scanner)
}

// ExtractSolutionData extracts data from pdb files and inserts it into a db.
func ExtractSolutionData(paths *bufio.Scanner) {
	// Construct the pipeline
	solutionGen, nLoaded := loadSolutions(paths)
	resultsGen := writeSolutions(solutionGen)
	// Wait for results
	for i := 0; i < nLoaded; i++ {
		<-resultsGen
	}
}

func loadSolutions(in *bufio.Scanner) (out chan *Solution, n int) {
	out = make(chan *Solution)
	for in.Scan() {
		go func(path string) {
			out <- NewSolution(path)
		}(in.Text())
		n++
	}
	return out, n
}

func writeSolutions(in chan *Solution) chan struct{} {
	out := make(chan struct{})
	for solution := range in {
		go func(s *Solution) {
			writeSolution(s)
			out <- struct{}{}
		}(solution)
	}
	return out
}

func writeSolution(solution *Solution) error {
	// Prepate the transaction with statements
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmtScores, err := prepareScoresTx(tx)
	if err != nil {
		return err
	}
	defer stmtScores.Close()
	err = solution.executeScores(stmtScores)
	if err != nil {
		return err
	}

	stmtActions, err := prepareActionsTx(tx)
	if err != nil {
		return err
	}
	defer stmtActions.Close()
	err = solution.executeActions(stmtActions)
	if err != nil {
		return err
	}

	stmtHistory, err := prepareHistoryTx(tx)
	if err != nil {
		return err
	}
	defer stmtHistory.Close()
	err = solution.executeHistory(stmtHistory)
	if err != nil {
		return err
	}

	stmtRank, err := prepareRankTx(tx)
	if err != nil {
		return err
	}
	defer stmtRank.Close()
	err = solution.executeRank(stmtRank)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
