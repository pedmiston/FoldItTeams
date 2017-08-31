/*foldit extracts data from FoldIt solution files.

Usage:
	foldit \
		-pdbs top_solution_file_paths.txt \
		-pdbType top \
		-outputDir top_solutions

Args:
	pdbs: Paths to solution files ending in a pdb extension. Paths can be provided
		in a file or via stdin.
	pdbType: Type of solution file. The two types of solution files are "regular"
		and "top". A top solution is just a regular solution with ranking data
		embedded in the name of the file.
	outputDir: Where to write the data. Since each solution generates multiple
		output files, the output must be a directory.
*/
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path"
)

var (
	pdbType, pdbs, outputDir *string
	scanner                  *bufio.Scanner
)

func main() {
	pdbs = flag.String("pdbs", "", "Files to process. Defaults to Stdin.")
	pdbType = flag.String("pdbType", "", "Type of solution file. Can be 'top' or 'regular'.")
	outputDir = flag.String("outputDir", "", "Destination for output files.")
	flag.Parse()

	var input *os.File
	var err error
	if *pdbs != "" {
		input, err = os.Open(*pdbs)
		if err != nil {
			log.Fatalln("Couldn't open input file")
		}
	} else {
		input = os.Stdin
	}
	scanner := bufio.NewScanner(input)

	switch *pdbType {
	case "top":
		ReadTopSolutions(scanner, *outputDir)
	case "regular":
		panic("Need to implement extractSolutionData")
	default:
		panic("unknown pdbType")
	}
}

func WriteTopSolutionData(scanner *bufio.Scanner, outputDir string) {
	topSolutionGen := LoadTopSolutions(scanner)
	tidyResultGen := TidyTopSolutions(topSolutionGen)

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

func loadTopSolutions(filenames *bufio.Scanner) <-chan *TopSolution {
	out := make(chan *TopSolution)
	go func() {
		for filenames.Scan() {
			out <- NewTopSolution(filenames.Text())
		}
		close(out)
	}()
	return out
}

func exportTopSolutions(in <-chan *TopSolution) <-chan *TopSolutionData {
	out := make(chan *TopSolutionData)
	go func() {
		for topSolution := range in {
			out <- NewTopSolutionData(topSolution)
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
