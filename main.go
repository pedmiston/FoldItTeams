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
	"flag"
	"log"
	"os"
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
	if *pdbs != "" {
		input, err := os.Open(*pdbs)
		if err != nil {
			log.Fatalln("Couldn't open input file")
		}
	} else {
		input := os.Stdin
	}
	scanner := bufio.NewScanner(input)

	switch *pdbType {
	case "top":
		extractTopSolutionData(scanner, *outputDir)
	case "regular":
		panic("Need to implement extractSolutionData")
	default:
		panic("unknown pdbType")
	}
}
