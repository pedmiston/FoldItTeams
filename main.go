/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -o=solution-data.json solution-filepaths.txt

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

	// Load solutions
	solutions, n := loadSolutions(scanner)

	for i := 0; i < n; i++ {
		s := <-solutions
		if err := encoder.Encode(s); err != nil {
			log.Println(err)
		}
	}
}
