/*foldit extracts data from FoldIt solution files.

Usage:
	foldit -input top_solution_file_paths.txt -type top -output top_solution_data.csv
	cat solution_file_paths.txt | foldit -type regular > solution_data.csv

File paths can be provided on stdin or in a file. Output gets sent
to the output file or stdout if none is given.

There are two types of solution files:

1. Regular solution files
2. Top solution files

Regular solution files are saved at regular intervals during gameplay. Top
solution files are solutions that have been evaluated for quality compared to all
other submissions. Top solution files contain data in the filename (e.g., the type
of ranking and the numeric rank of this particular solution.).
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	inputFile, outputFile *os.File
	err                   error
)

type result struct {
	data TopSolution
	err  error
}

func main() {
	solutionType := flag.String("type", "", "Type of solution file. Can be 'top' or 'regular'.")
	input := flag.String("input", "", "Files to process. Defaults to Stdin.")
	output := flag.String("output", "", "Destination for data. Defaults to Stdout.")
	flag.Parse()

	if *input != "" {
		inputFile, err = os.Open(*input)
		if err != nil {
			log.Fatalln("Problem opening input file")
		}
	} else {
		inputFile = os.Stdin
	}

	if *output != "" {
		outputFile, err = os.Create(*output)
		if err != nil {
			log.Fatalln("Error opening output file:", *output)
		}
	} else {
		outputFile = os.Stdout
	}

	switch *solutionType {
	case "top":
		processTopSolutions(inputFile, outputFile)
	case "regular":
		panic("Need to implement parseRegularSolution")
	default:
		panic("unknown solutionType")
	}

}

func processTopSolutions(input *os.File, output *os.File) {
	scanner := bufio.NewScanner(input)
	encoder := json.NewEncoder(output)

	// Run a go routine for each input file
	// and send the results back on a channel.
	ch := make(chan result)
	var chSize int
	for scanner.Scan() {
		go func(f string) {
			topSolution, err := readTopSolution(f)
			ch <- result{*topSolution, err}
		}(scanner.Text())
		chSize++
	}

	// Pull results from the channel.
	for j := 0; j < chSize; j++ {
		result := <-ch
		if result.err != nil {
			log.Println(result.err)
		}
		err := encoder.Encode(&result.data)
		if err != nil {
			log.Println(err)
		}
	}
}
