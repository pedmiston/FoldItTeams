package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	scoresFields  = []string{"filename", "puzzle_id", "user_id", "group_id", "score"}
	actionsFields = []string{"filename", "action", "count"}
	historyFields = []string{"filename", "ix", "id"}
	rePuzzleID    = regexp.MustCompile(`solution_(?P<PuzzleID>\d+)/`)
	rePDL         = regexp.MustCompile(`^IRDATA PDL`)
	reTimestamp   = regexp.MustCompile(`^IRDATA TIMESTAMP`)
	reHistory     = regexp.MustCompile(`^IRDATA HISTORY`)
)

// A Solution is a collection of data extracted from a solution file.
type Solution struct {
	PuzzleID  int
	UserID    int
	GroupID   int
	Score     float64
	Timestamp int
	Actions   map[string]int
	History   []string
	Filename  string
	Errors    []error
}

// NewSolution creates a new Solution from solution pdb file.
func NewSolution(filename string) *Solution {
	// The minimum Solution contains only the solution filename
	solution := &Solution{Filename: filename}

	file, err := os.Open(filename)
	if err != nil {
		solution.Errors = []error{err}
	}
	defer file.Close()

	puzzleID, err := readPuzzleIDFromFilename(filename)
	if err != nil {
		solution.Errors = append(solution.Errors, err)
	}
	solution.PuzzleID = puzzleID

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// TODO
		if rePDL.MatchString(line) {
			solution.addDataFromPDL(line)
		} else if reTimestamp.MatchString(line) {
			solution.Timestamp, _ = strconv.Atoi(strings.Split(line, " ")[2:][0])
		} else if reHistory.MatchString(line) {
			solution.History = strings.Split(strings.Split(line, " ")[2:][0], ",")
		}
	}

	return solution
}

func (s *Solution) addDataFromPDL(pdlLine string) {
	split := strings.Split(pdlLine, ",")
	s.UserID, _ = strconv.Atoi(split[2])
	s.GroupID, _ = strconv.Atoi(split[3])

	lastItem := strings.Split(split[7], " ")
	s.Score, _ = strconv.ParseFloat(lastItem[0], 64)

	if len(lastItem) > 3 {
		s.addActions(lastItem[3:])
	}
}

func (s *Solution) addActions(actions []string) {
	if s.Actions == nil {
		s.Actions = make(map[string]int)
	}

	for _, value := range actions {
		prefix := strings.Split(value, "|")
		keyValue := prefix[len(prefix)-1]
		items := strings.Split(keyValue, "=")
		if len(items) == 2 && items[0] != "" {
			count, _ := strconv.Atoi(items[1])
			prev, exists := s.Actions[items[0]]
			if exists {
				s.Actions[items[0]] = prev + count
			} else {
				s.Actions[items[0]] = count
			}
		}
	}
}

func (s *Solution) prepare(tx *sql.Tx, tblName string, fields []string) *sql.Stmt {
	q := make([]string, len(fields))
	for i := range fields {
		q[i] = "?"
	}

	msg := fmt.Sprintf("insert into %s(%s) values(%s)",
		tblName, strings.Join(fields, ","), strings.Join(q, ","))

	stmt, err := tx.Prepare(msg)
	if err != nil {
		log.Fatalf("%s:\n%s", err, msg)
	}

	return stmt
}

func (s *Solution) prepareScores(tx *sql.Tx) *sql.Stmt {
	return s.prepare(tx, "scores", scoresFields)
}

func (s *Solution) prepareActions(tx *sql.Tx) *sql.Stmt {
	return s.prepare(tx, "actions", actionsFields)
}

func (s *Solution) prepareHistory(tx *sql.Tx) *sql.Stmt {
	return s.prepare(tx, "history", historyFields)
}

func (s *Solution) getScores() (values []string) {
	values = []string{
		s.Filename,
		strconv.Itoa(s.PuzzleID),
		strconv.Itoa(s.UserID),
		strconv.Itoa(s.GroupID),
		strconv.FormatFloat(s.Score, 'f', -1, 64),
	}
	return
}

func (s *Solution) getActions() [][]string {
	records := make([][]string, len(s.Actions))
	var row int
	for action, count := range s.Actions {
		records[row] = []string{
			s.Filename,
			action,
			strconv.Itoa(count),
		}
		row++
	}
	return records
}

func (s *Solution) getHistory() [][]string {
	records := make([][]string, len(s.History))
	for ix, id := range s.History {
		records[ix] = []string{
			s.Filename,
			strconv.Itoa(ix),
			id,
		}
	}
	return records
}

func readPuzzleIDFromFilename(solutionFilename string) (puzzleID int, err error) {
	puzzleIDMatch := rePuzzleID.FindStringSubmatch(solutionFilename)
	if len(puzzleIDMatch) != 2 {
		return 0, errors.New("rePuzzleID not matched by " + solutionFilename)
	}
	// Ignoring conversion error because the regexp only matches ints
	puzzleID, _ = strconv.Atoi(puzzleIDMatch[1])
	return puzzleID, nil
}
