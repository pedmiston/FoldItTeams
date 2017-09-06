package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	rePuzzleID   = regexp.MustCompile(`solution_(?P<PuzzleID>\d+)/`)
	rePDL        = regexp.MustCompile(`^IRDATA PDL`)
	reTimestamp  = regexp.MustCompile(`^IRDATA TIMESTAMP`)
	reHistory    = regexp.MustCompile(`^IRDATA HISTORY`)
	scoresFields = []string{"filename", "puzzle_id", "user_id",
		"group_id", "score"}
	actionsFields = []string{"filename", "action", "count"}
	historyFields = []string{"filename", "ix", "id"}
	rankFields    = []string{"filename", "type", "rank"}
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
	RankType  string
	Rank      int
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

	rankType, rank, err := getRankFromFilename(filename)
	if err != nil {
		solution.Errors = append(solution.Errors, err)
	}
	solution.RankType = rankType
	solution.Rank = rank

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

func prepareTx(tx *sql.Tx, tblName string, fields []string) (*sql.Stmt, error) {
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

func prepareScoresTx(tx *sql.Tx) (*sql.Stmt, error) {
	return prepareTx(tx, "scores", scoresFields)
}

func prepareActionsTx(tx *sql.Tx) (*sql.Stmt, error) {
	return prepareTx(tx, "actions", actionsFields)
}

func prepareHistoryTx(tx *sql.Tx) (*sql.Stmt, error) {
	return prepareTx(tx, "history", historyFields)
}

func prepareRankTx(tx *sql.Tx) (*sql.Stmt, error) {
	return prepareTx(tx, "ranks", rankFields)
}

func (s *Solution) executeScores(stmt *sql.Stmt) error {
	_, err := stmt.Exec(
		s.Filename,
		s.PuzzleID,
		s.UserID,
		s.GroupID,
		s.Score,
	)
	return err
}

func (s *Solution) executeActions(stmt *sql.Stmt) error {
	for action, count := range s.Actions {
		_, err := stmt.Exec(s.Filename, action, count)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Solution) executeHistory(stmt *sql.Stmt) error {
	for ix, id := range s.History {
		_, err := stmt.Exec(s.Filename, ix, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Solution) executeRank(stmt *sql.Stmt) error {
	if s.RankType != "" {
		_, err := stmt.Exec(
			s.Filename,
			s.PuzzleID,
			s.UserID,
			s.GroupID,
			s.Score,
		)
		return err
	}
	return nil
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
