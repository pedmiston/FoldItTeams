/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -t=json -o=data.json filepaths.txt
*/
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	format := flag.String("t", "", "output type")
	output := flag.String("o", "", "output file")
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

	switch *format {
	case "json":
		Write(src, dst)
	default:
		panic("unknown output type")
	}
}
