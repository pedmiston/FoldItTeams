package main

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	reRankInfo = regexp.MustCompile(
		`solution_(?P<RankType>[a-z]+)_(?P<Rank>\d+)_\d+_\d+_\d+.ir_solution.pdb`)
)

// A TopSolution is a Solution with ranking information
type TopSolution struct {
	Solution
	RankType string
	Rank     int
}

// readTopSolution extracts a record of performance data from
// a top solution file.
func readTopSolution(name string) (topSolution *TopSolution, err error) {
	solution, err := readSolution(name)
	rankType, rank, _ := readRankFromFilename(name)
	topSolution = &TopSolution{
		Solution: solution,
		RankType: rankType,
		Rank:     rank,
	}
	return topSolution, err
}

func readRankFromFilename(name string) (rankType string, rank int, err error) {
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
