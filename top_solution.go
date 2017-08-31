package main

import (
	"encoding/csv"
	"errors"
	"regexp"
	"strconv"
)

// A TopSolution is a solution with ranking information
type TopSolution struct {
	*Solution
	RankType string
	Rank     int
}

var (
	reRankInfo = regexp.MustCompile(
		`solution_(?P<RankType>[a-z]+)_(?P<Rank>\d+)_\d+_\d+_\d+.ir_solution.pdb`)
)

// NewTopSolution creates a new TopSolution from a top solution pdb file.
func NewTopSolution(name string) *TopSolution {
	solution := NewSolution(name)
	rankType, rank, err := getRankFromFilename(name)

	if err != nil {
		solution.Errors = append(solution.Errors, err)
	}

	topSolution := &TopSolution{
		Solution: solution,
		RankType: rankType,
		Rank:     rank,
	}

	return topSolution
}

func (t *TopSolution) writeScores(writer *csv.Writer) {
	data := []string{
		strconv.Itoa(t.PuzzleID),
		strconv.Itoa(t.UserID),
		strconv.Itoa(t.GroupID),
		strconv.FormatFloat(t.Score, 'f', -1, 64),
		t.RankType,
		strconv.Itoa(t.Rank),
		t.Filename,
	}
	writer.Write(data)
}

func (t *TopSolution) writeActions(writer *csv.Writer) {
	var data []string
	for action, count := range t.Actions {
		data = []string{
			t.Filename,
			action,
			strconv.Itoa(count),
		}
		writer.Write(data)
	}
}

func (t *TopSolution) writeHistory(writer *csv.Writer) {
	var data []string
	for ix, id := range t.History {
		data = []string{
			t.Filename,
			strconv.Itoa(ix),
			id,
		}
		writer.Write(data)
	}
}

// getRankFromFilename extracts rank and rank type from a solution filename.
func getRankFromFilename(name string) (rankType string, rank int, err error) {
	matches := reRankInfo.FindAllStringSubmatch(name, -1)
	if len(matches) == 0 {
		err = errors.New("Unable to read rank info from filename: " + name)
		return
	}
	matchValues := matches[0]
	rankType = matchValues[1]
	rank, err = strconv.Atoi(matchValues[2])
	if err != nil {
		err = errors.New("Unable to convert rank to integer: " + matchValues[2])
	}
	return
}
