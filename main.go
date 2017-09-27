/*foldit-scraper extracts IRDATA data fields from FoldIt solution files.

Usage:
	scraper filepaths.txt > data.json
	find . -name "*.pdb" | scraper > data.json
*/
package main

import (
	"log"
	"os"
)

func main() {
	var src = os.Stdin
	var err error

	if len(os.Args) == 2 {
		input := os.Args[1]
		src, err = os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
	}

	Scrape(src, os.Stdout, os.Stderr)
}
