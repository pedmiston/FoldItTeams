package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
)

// A TopSolution is a solution with ranking information
type TopSolution struct {
	Solution
	RankType string
	Rank     int
}

var (
	reRankInfo = regexp.MustCompile(
		`solution_(?P<RankType>[a-z]+)_(?P<Rank>\d+)_\d+_\d+_\d+.ir_solution.pdb`)
)

type tidyResult struct {
	Scores  [][]string
	Actions [][]string
	History [][]string
}

// ReadTopSolutions reads solution files from a scanner
// and collects the results in the outputDir.
func ReadTopSolutions(scanner *bufio.Scanner, outputDir string) {
	topSolutionGen := readAllTopSolutions(scanner)
	tidyResultGen := tidyAllTopSolutions(topSolutionGen)

	scoresWriter := createWriter(outputDir, "scores")
	actionsWriter := createWriter(outputDir, "actions")
	historyWriter := createWriter(outputDir, "history")

	for tidyResult := range tidyResultGen {
		for _, line := range tidyResult.Scores {
			scoresWriter.Write(line)
		}
		for _, line := range tidyResult.Actions {
			actionsWriter.Write(line)
		}
		for _, line := range tidyResult.History {
			historyWriter.Write(line)
		}
	}
}

// readAllTopSolutions loads a channel with TopSolutions
// created from solution files.
func readAllTopSolutions(scanner *bufio.Scanner) <-chan *TopSolution {
	out := make(chan *TopSolution)
	go func() {
		for scanner.Scan() {
			out <- readTopSolution(scanner.Text())
		}
		close(out)
	}()
	return out
}

// readTopSolution creates a new TopSolution from a top solution file.
func readTopSolution(name string) *TopSolution {
	solution, _ := readSolution(name)
	rankType, rank, _ := getRankFromFilename(name)
	topSolution := &TopSolution{
		Solution: solution,
		RankType: rankType,
		Rank:     rank,
	}
	return topSolution
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

func tidyAllTopSolutions(in <-chan *TopSolution) <-chan *tidyResult {
	out := make(chan *tidyResult)
	go func() {
		for topSolution := range in {
			out <- tidyTopSolution(topSolution)
		}
		close(out)
	}()
	return out
}

func tidyTopSolution(topSolution *TopSolution) *tidyResult {
	tidyResult := &tidyResult{
		Scores:  topSolution.getScores(),
		Actions: topSolution.getActions(),
		History: topSolution.getHistory(),
	}
	return tidyResult
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

func createWriter(outputDir, dataType string) *csv.Writer {
	outputDst := path.Join(outputDir, dataType+".csv")
	outputFile, err := os.Open(outputDst)
	if err != nil {
		log.Fatalln("Can't open output file " + outputDst)
	}
	outputWriter := csv.NewWriter(outputFile)
	return outputWriter
}
