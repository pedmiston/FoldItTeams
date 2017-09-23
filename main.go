/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -t=json filepaths.txt > data.json
	foldit -t=mysql data.json
	foldit -t=json filepaths.txt | foldit -t=mysql
*/
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	format := flag.String("t", "", "output type")
	flag.Parse()

	// Open the input file
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

	switch *format {
	case "json":
		Write(src, os.Stdout)
	default:
		panic("unknown output type")
	}
}
