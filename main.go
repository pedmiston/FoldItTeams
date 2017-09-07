/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -paths=solution-filepaths.txt -dest=foldit.sqlite

Args:
	paths: A file containing paths to solution files. If a file is not provided,
		paths are expected via stdin.
	dest: Location of a sqlite database containing the extracted data.
*/
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	paths := flag.String("paths", "",
		"A file containing paths to files to process.")
	dest := flag.String("dest", "", "Location of a sqlite database")
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
	if *dest == "" {
		*dest = "foldit.db"
	}
	os.Remove(*dest)
	db, err = sql.Open("sqlite3", *dest)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tables in the db
	err = createTables()
	if err != nil {
		log.Fatal(err)
	}

	// Load solutions concurrently
	solutionGen, nLoaded := loadSolutions(scanner)
	// Write solutions concurrently
	finished := writeSolutions(solutionGen, nLoaded)
	// Wait for all solutions to finish
	for i := 0; i < nLoaded; i++ {
		<-finished
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

func writeSolutions(in chan *Solution, n int) (out chan bool) {
	out = make(chan bool)
	for i := 0; i < n; i++ {
		solution := <-in
		go func(s *Solution) {
			err := writeSolution(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "writing solution: %v\n%v", err, s.Filename)
			}
			out <- true
		}(solution)
	}
	return out
}

func writeSolution(solution *Solution) error {
	// Prepare the transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the statements in this transaction
	// TODO Refactor!!!
	stmtScores, err := prepareStatement(tx, "scores", scoresFields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "preparing scores statement: %v", err)
	} else {
		defer stmtScores.Close()
		err = solution.executeScores(stmtScores)
		if err != nil {
			fmt.Fprintf(os.Stderr, "executing scores statement: %v", err)
		}
	}

	stmtActions, err := prepareStatement(tx, "actions", actionsFields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "preparting actions statement: %v", err)
	} else {
		defer stmtActions.Close()
		err = solution.executeActions(stmtActions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "executing actions statement: %v", err)
		}
	}

	stmtHistory, err := prepareStatement(tx, "history", historyFields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "preparing history statement: %v", err)
	} else {
		defer stmtHistory.Close()
		err = solution.executeHistory(stmtHistory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "executing history statement: %v", err)
		}
	}

	stmtRank, err := prepareStatement(tx, "ranks", rankFields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "preparing rank statement: %v", err)
	} else {
		defer stmtRank.Close()
		err = solution.executeRank(stmtRank)
		if err != nil {
			fmt.Fprintf(os.Stderr, "executing rank statement: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
