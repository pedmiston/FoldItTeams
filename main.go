/*foldit extracts data from FoldIt solution files.

Usage:
	foldit filepaths.txt > data.json
*/
package main

import (
	"log"
	"os"
)

func main() {
	// Open the input file
	var src *os.File
	var err error
	if len(os.Args) == 1 {
		src = os.Stdin
	} else {
		input := os.Args[1]
		src, err = os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer src.Close()

	Write(src, os.Stdout)
}
