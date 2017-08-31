package main

import (
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

type TopSolutionData struct {
	Scores  [][]string
	Actions [][]string
	History [][]string
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

func (t *TopSolution) getScores() [][]string {
	return [][]string{[]string{
		strconv.Itoa(t.PuzzleID),
		strconv.Itoa(t.UserID),
		strconv.Itoa(t.GroupID),
		strconv.FormatFloat(t.Score, 'f', -1, 64),
		t.Filename,
	}}
}

func (t *TopSolution) getActions() [][]string {
	var row int
	dataframe := make([][]string, len(t.Actions))
	for action, count := range t.Actions {
		dataframe[row] = []string{
			strconv.Itoa(t.PuzzleID),
			strconv.Itoa(t.UserID),
			strconv.Itoa(t.GroupID),
			strconv.FormatFloat(t.Score, 'f', -1, 64),
			t.Filename,
			action,
			strconv.Itoa(count),
		}
		row++
	}
	return dataframe
}

func (t *TopSolution) getHistory() [][]string {
	return [][]string{[]string{
		strconv.Itoa(t.PuzzleID),
		strconv.Itoa(t.UserID),
		strconv.Itoa(t.GroupID),
		strconv.FormatFloat(t.Score, 'f', -1, 64),
		t.Filename,
	}}
}

func NewTopSolutionData(topSolution *TopSolution) *TopSolutionData {
	topSolutionData := &TopSolutionData{
		Scores:  topSolution.getScores(),
		Actions: topSolution.getActions(),
		History: topSolution.getHistory(),
	}

	return topSolutionData
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
