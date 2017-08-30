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

func extractTopSolutionData(scanner *bufio.Scanner, outputDir string) {
	topSolutioner := readAllTopSolutionFiles(scanner)
	tidySolutioner := tidyAllTopSolutions(topSolutioner)

	scoresWriter := createWriter(outputDir, "scores")
	actionsWriter := createWriter(outputDir, "actions")
	historyWriter := createWriter(outputDir, "history")

	for tidyResult := range tidySolutioner {
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

func readAllTopSolutionFiles(scanner *bufio.Scanner) <-chan *TopSolution {
	out := make(chan *TopSolution)
	go func() {
		for scanner.Scan() {
			out <- readTopSolutionFile(scanner.Text())
		}
		close(out)
	}()
	return out
}

func readTopSolutionFile(name string) *TopSolution {
	solution, err := readSolution(name)
	rankType, rank, _ := readRankFromFilename(name)
	topSolution := &TopSolution{
		Solution: solution,
		RankType: rankType,
		Rank:     rank,
	}
	return topSolution
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

func tidyAllTopSolutions(in <-chan *TopSolution) <-chan tidyResult {
	out := make(chan tidyResult)
	go func() {
		for topSolution := range in {
			out <- tidyTopSolution(topSolution)
		}
		close(out)
	}()
	return out
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
