/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -i=solution-filepaths.txt -o=solution-data.json

Args:
	paths: A file containing paths to solution files. If a file is not provided,
		paths are expected via stdin.
	dest: Location of a sqlite database containing the extracted data.
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
)

func main() {
	input := flag.String("i", "", "input file, e.g., filepaths.txt")
	output := flag.String("o", "solutions.json", "output file")

	// Create a solution file path scanner
	var src *os.File
	var err error
	if *input == "" {
		src = os.Stdin
	} else {
		src, err = os.Open(*input)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer src.Close()
	scanner := bufio.NewScanner(src)

	// Create the output file
	dst, err := os.Create(*output)
	defer dst.Close()
	encoder := json.NewEncoder(dst)

	// Load solutions
	solutions, n := loadSolutions(scanner)

	for i := 0; i < n; i++ {
		s := <-solutions
		if err := encoder.Encode(s); err != nil {
			log.Println(err)
		}
	}
}
