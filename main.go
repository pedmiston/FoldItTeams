/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -o=data.json filepaths.txt
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
	output := flag.String("o", "", "output file")
	flag.Parse()

	// Create a solution file path scanner
	var src *os.File
	var err error
	if len(flag.Args()) == 0 {
		src = os.Stdin
	} else {
		input := flag.Args()[0]
		src, err = os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer src.Close()
	scanner := bufio.NewScanner(src)

	// Load solutions
	solutions, n := loadSolutions(scanner)

	// Create the output file
	var dst *os.File
	if *output == "" {
		dst = os.Stdout
	} else {
		dst, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer dst.Close()
	encoder := json.NewEncoder(dst)

	// Write solution data
	writeSolutions(solutions, n, encoder)
}
