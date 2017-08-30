/*foldit extracts data from FoldIt solution files.

Subcommands:
	foldit new [options]   # convert new solution files to json data
	foldit tidy [options]  # convert json data to structured csvs


# foldit new [options]

The `new` subcommand processes new solution files in batches.

Usage:
	foldit new -input top_solution_file_paths.txt -type top -output top_solution_data.json
	cat top_solution_file_paths.txt | foldit new -type top > top_solution_data.json

File paths can be provided on stdin or in a file. Output gets sent
to the output file or stdout if none is given.

There are two types of solution files:

1. Regular solution files
2. Top solution files

Regular solution files are saved at regular intervals during gameplay. Top
solution files are solutions that have been evaluated for quality compared to all
other submissions. Top solution files contain data in the filename (e.g., the type
of ranking and the numeric rank of this particular solution.).


# foldit tidy [options]

The `tidy` subcommand assembles csvs from all batched solution files. There
are three types of tidy outputs that can be created from each solution struct.

1. Scores
2. Actions
3. History

Usage:
	foldit tidy -input top_solution_data.json -type actions -output top_solution_actions.csv
*/
package main

import (
	"flag"
	"log"
	"os"
)

var (
	solutionType, input, output *string
	inputFile, outputFile       *os.File
	err                         error
)

func main() {
	solutionType = flag.String("type", "", "Type of solution file. Can be 'top' or 'regular'.")
	input = flag.String("input", "", "Files to process. Defaults to Stdin.")
	output = flag.String("output", "", "Destination for data. Defaults to Stdout.")
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

	subcommand := flag.Args()[0]
	switch subcommand {
	case "new":
		newSolutions(inputFile, outputFile)
	case "tidy":
		tidy(inputFile, outputFile)
	default:
		log.Fatalln("first arg " + subcommand + " must be 'new' or 'tidy'")
	}
}

func newSolutions(input *os.File, output *os.File) {
	switch *solutionType {
	case "top":
		writeTopSolutionsToJSON(inputFile, outputFile)
	case "regular":
		panic("Need to implement parseRegularSolution")
	default:
		panic("unknown solutionType")
	}
}

func tidy(input *os.File, output *os.File) {
	switch *solutionType {
	case "scores":
		//
	case "actions":
		//
	default:
		panic("unknown solutionType")
	}
}
