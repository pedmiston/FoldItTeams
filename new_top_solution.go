/*Create TopSolution structs from data in PDB files*/
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	reRankInfo = regexp.MustCompile(
		`solution_(?P<RankType>[a-z]+)_(?P<Rank>\d+)_\d+_\d+_\d+.ir_solution.pdb`)
)

// A TopPDBResult is the result of reading a top solution pdb file
type TopPDBResult struct {
	data TopSolution
	err  error
}

// writeTopSolutionsToJSON extracts data from top solution pdb files that
// have been externally ranked and encodes it in JSON format to the output.
func writeTopSolutionsToJSON(input *os.File, output *os.File) {
	scanner := bufio.NewScanner(input)
	encoder := json.NewEncoder(output)

	// Run a go routine for each input file
	// and send the results back on a channel.
	ch := make(chan TopPDBResult)
	var chSize int
	for scanner.Scan() {
		go func(f string) {
			topSolution, err := readTopSolutionFromPDB(f)
			ch <- TopPDBResult{*topSolution, err}
		}(scanner.Text())
		chSize++
	}

	// Pull results from the channel.
	for j := 0; j < chSize; j++ {
		result := <-ch
		if result.err != nil {
			log.Println(result.err)
		}
		err := encoder.Encode(&result.data)
		if err != nil {
			log.Println(err)
		}
	}
}

// readTopSolutionFromPDB extracts a record of performance data from
// a top solution pdb file.
func readTopSolutionFromPDB(name string) (topSolution *TopSolution, err error) {
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
