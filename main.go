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
		WriteTopData(scanner, *outputDir)
	case "regular":
		panic("Need to implement extractSolutionData")
	default:
		panic("unknown pdbType")
	}
}

func WriteTopData(scanner *bufio.Scanner, outputDir string) {
	genTopSolution := loadTopSolutions(scanner)
	genTopData := pushTopData(genTopSolution)

	scoresWriter := newWriter(outputDir, "scores")
	actionsWriter := newWriter(outputDir, "actions")
	historyWriter := newWriter(outputDir, "history")

	for topData := range genTopData {
		for _, record := range topData.Scores {
			scoresWriter.Write(record)
		}
		for _, record := range topData.Actions {
			actionsWriter.Write(record)
		}
		for _, record := range topData.History {
			historyWriter.Write(record)
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

func pushTopData(in <-chan *TopSolution) <-chan *TopData {
	out := make(chan *TopData)
	go func() {
		for topSolution := range in {
			out <- NewTopData(topSolution)
		}
		close(out)
	}()
	return out
}

func newWriter(outputDir, dataType string) *csv.Writer {
	outputDst := path.Join(outputDir, dataType+".csv")
	outputFile, err := os.Open(outputDst)
	if err != nil {
		log.Fatalln("Can't open output file " + outputDst)
	}
	outputWriter := csv.NewWriter(outputFile)
	return outputWriter
}
