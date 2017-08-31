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
	if *pdbs == "" {
		input = os.Stdin
	} else {
		input, err = os.Open(*pdbs)
		if err != nil {
			log.Fatalln("Couldn't open input file")
		}
	}
	scanner := bufio.NewScanner(input)

	switch *pdbType {
	case "top":
		WriteTopData(scanner, *outputDir)
	case "regular":
		panic("Need to implement WriteRegularData")
	default:
		panic("unknown pdbType")
	}
}

// WriteTopData writes data from top solution pdb files to the output dir.
func WriteTopData(topSolutionFilenames *bufio.Scanner, outputDir string) {
	genTopSolution := loadTopSolutions(topSolutionFilenames)

	scoresWriter := newWriter(outputDir, "scores.csv")
	actionsWriter := newWriter(outputDir, "actions.csv")
	historyWriter := newWriter(outputDir, "history.csv")

	for topSolution := range genTopSolution {
		topSolution.writeScores(scoresWriter)
		topSolution.writeActions(actionsWriter)
		topSolution.writeHistory(historyWriter)
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

func newWriter(outputDir, dataType string) *csv.Writer {
	outputDst := path.Join(outputDir, dataType)
	outputFile, err := os.Open(outputDst)
	if err != nil {
		log.Fatalln("Can't open output file " + outputDst)
	}
	outputWriter := csv.NewWriter(outputFile)
	return outputWriter
}
