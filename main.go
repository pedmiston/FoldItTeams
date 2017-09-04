/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -pdbs=paths.txt -pdbType=top -outputDir=top

Args:
	pdbs: A file containing paths to pdb files. If a file is not provided,
		paths are expected via stdin.
	pdbType: Type of solution file. The two types of solution files are
		"regular" and "top". A top solution is just a regular solution
		with ranking data embedded in the filename.
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
	pdbs, pdbType, outputDir *string
	scanner                  *bufio.Scanner
)

func main() {
	pdbs = flag.String("pdbs", "",
		"A file containing paths to files to process.")
	pdbType = flag.String("pdbType", "",
		"Type of solution file. Can be 'top' or 'regular'.")
	outputDir = flag.String("outputDir", "",
		"Destination for output files.")
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

// WriteTopData writes data from top ranked solutions to the output dir.
func WriteTopData(filenames *bufio.Scanner, outputDir string) {
	genTopSolution := loadTopSolutions(filenames)

	scoresFile := createOutputFile(outputDir, "scores.csv")
	defer scoresFile.Close()

	actionsFile := createOutputFile(outputDir, "actions.csv")
	defer actionsFile.Close()

	historyFile := createOutputFile(outputDir, "history.csv")
	defer historyFile.Close()

	scoresWriter, actionsWriter, historyWriter := createWriters(
		scoresFile, actionsFile, historyFile)

	// Pull topSolutions out of the chan and write each one
	for topSolution := range genTopSolution {
		topSolution.writeScoresTo(scoresWriter)
		topSolution.writeActionsTo(actionsWriter)
		topSolution.writeHistoryTo(historyWriter)
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

func createOutputFile(outputDir, filename string) *os.File {
	filePath := path.Join(outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Unable to create output file %s: %s\n", filePath, err)
	}
	return file
}

func createWriters(writers ...*os.File) []*csv.Writer {
	var csvWriters = []*csv.Writer{}
	for i, f := range writers {
		csvWriters[i] = csv.NewWriter(f)
	}
	return csvWriters
}
