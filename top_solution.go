package main

import (
	"errors"
	"regexp"
	"strconv"
)

var reRankInfo = regexp.MustCompile(
	`solution_(?P<RankType>[a-z]+)_(?P<Rank>\d+)_\d+_\d+_\d+.ir_solution.pdb`)

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
