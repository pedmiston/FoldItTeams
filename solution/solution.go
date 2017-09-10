package solution

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reSID          = regexp.MustCompile(`^IRDATA SID (\d+)`)
	rePID          = regexp.MustCompile(`^IRDATA PID (\d+)`)
	reScore        = regexp.MustCompile(`^IRDATA SCORE (\d+\.?\d*)`)
	reTimestamp    = regexp.MustCompile(`^IRDATA TIMESTAMP (\d+)`)
	reHistory      = regexp.MustCompile(`^IRDATA HISTORY (.*)`)
	reMacroHistory = regexp.MustCompile(`^IRDATA MACRO_HIST (.*)`)
	reMoves        = regexp.MustCompile(`^IRDATA SOLN_MOVECOUNT (\d+)`)
	rePDL          = regexp.MustCompile(`^IRDATA PDL \.+ (.*)`)
)

// A Solution is a collection of data extracted from a solution file.
type Solution struct {
	SolutionID   int
	PuzzleID     int
	Timestamp    time.Time
	UserID       int
	GroupID      int
	Actions      map[string]int
	Score        float64
	History      []string
	MacroHistory []string
	Filepath     string
	Moves        int
}

// New creates a new Solution from solution pdb file.
func New(filename string) (s *Solution, err error) {
	s = &Solution{Filepath: filename}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Scan the lines in the solution file,
	// adding data to the solution.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case reSID.MatchString(line):
			err = s.extractSolutionID(line)
		case rePID.MatchString(line):
			err = s.extractPuzzleID(line)
		case reScore.MatchString(line):
			err = s.extractScore(line)
		case reTimestamp.MatchString(line):
			err = s.extractTimestamp(line)
		case rePDL.MatchString(line):
			err = s.addDataFromPDL(line)
		case reHistory.MatchString(line):
			err = s.extractHistory(line)
		case reMacroHistory.MatchString(line):
			err = s.extractMacroHistory(line)
		case reMoves.MatchString(line):
			err = s.extractMoves(line)
		}

		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Solution) extractSolutionID(line string) (err error) {
	s.SolutionID, err = strconv.Atoi(reSID.FindStringSubmatch(line)[1])
	return err
}

func (s *Solution) extractPuzzleID(line string) (err error) {
	s.PuzzleID, err = strconv.Atoi(rePID.FindStringSubmatch(line)[1])
	return err
}

func (s *Solution) extractScore(line string) (err error) {
	s.Score, err = strconv.ParseFloat(reScore.FindStringSubmatch(line)[1], 64)
	return err
}

func (s *Solution) extractTimestamp(line string) (err error) {
	timeint, err := strconv.ParseInt(reTimestamp.FindStringSubmatch(line)[1], 10, 64)
	if err != nil {
		return err
	}
	s.Timestamp = time.Unix(timeint, 0)
	return
}

func (s *Solution) extractHistory(line string) (err error) {
	if s.History != nil {
		return errors.New("attempting to overwrite s.History")
	}
	s.History = strings.Split(reHistory.FindStringSubmatch(line)[1], ",")
	return
}

func (s *Solution) extractMacroHistory(line string) (err error) {
	if s.MacroHistory != nil {
		return errors.New("attempting to overwriter s.MacroHistory")
	}
	s.MacroHistory = strings.Split(reMacroHistory.FindStringSubmatch(line)[1], ",")
	return
}

func (s *Solution) extractMoves(line string) (err error) {
	s.Moves, err = strconv.Atoi(reMoves.FindStringSubmatch(line)[1])
	return err
}

func (s *Solution) addDataFromPDL(line string) (err error) {
	split := strings.Split(line, ",")

	userID, err := strconv.Atoi(split[2])
	if err != nil {
		return err
	}
	if s.UserID != 0 && s.UserID != userID {
		return errors.New("attempting to overwrite s.UserID")
	}
	s.UserID = userID

	groupID, err := strconv.Atoi(split[3])
	if err != nil {
		return err
	}
	if s.GroupID != 0 && s.GroupID != groupID {
		return errors.New("attempting to overwrite s.GroupID")
	}
	s.GroupID = groupID

	lastItem := strings.Split(split[7], " ")
	if len(lastItem) > 3 {
		err = s.addActions(lastItem[2:])
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *Solution) addActions(actions []string) (err error) {
	if s.Actions == nil {
		s.Actions = make(map[string]int)
	}

	for _, value := range actions {
		prefix := strings.Split(value, "|")
		keyValue := prefix[len(prefix)-1]
		items := strings.Split(keyValue, "=")
		if len(items) == 2 && items[0] != "" {
			count, err := strconv.Atoi(items[1])
			if err != nil {
				return err
			}
			s.Actions[items[0]] += count
		}
	}

	return nil
}
