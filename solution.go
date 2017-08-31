package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	rePuzzleID  = regexp.MustCompile(`solution_(?P<PuzzleID>\d+)/`)
	rePDL       = regexp.MustCompile(`^IRDATA PDL`)
	reTimestamp = regexp.MustCompile(`^IRDATA TIMESTAMP`)
	reHistory   = regexp.MustCompile(`^IRDATA HISTORY`)
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

func readPuzzleIDFromFilename(solutionFilename string) (puzzleID int, err error) {
	puzzleIDMatch := rePuzzleID.FindStringSubmatch(solutionFilename)
	if len(puzzleIDMatch) != 2 {
		return 0, errors.New("rePuzzleID not matched by " + solutionFilename)
	}
	// Ignoring conversion error because the regexp only matches ints
	puzzleID, _ = strconv.Atoi(puzzleIDMatch[1])
	return puzzleID, nil
}
