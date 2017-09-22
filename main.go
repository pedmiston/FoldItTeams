/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -o=data.json filepaths.txt
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pedmiston/foldit/irdata"
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

	// Load a channel of IRDATA from the input scanner.
	out, n := irdata.Load(scanner)

	// Pull data from the channel and encode it to JSON.
	for i := 0; i < n; i++ {
		r := <-out
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "%s,%v\n", r.Data["Filepath"], r.Err)
			continue
		}
		if err := encoder.Encode(r.Data); err != nil {
			fmt.Fprintf(os.Stderr, "%s,%v\n", r.Data["Filepath"], r.Err)
		}
	}

	// Write solution data to the output encoder
	writeSolutions(solutions, n, encoder)
}
