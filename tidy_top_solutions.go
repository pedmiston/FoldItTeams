/*Reading JSON files and returning slices of TopSolution structs.*/
package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// A TopSolutionsResult is the result of reading top solutions from a json file
type TopSolutionsResult struct {
	solutions []TopSolution
}

// tidyTopSolutions
func tidyTopSolutions(input *os.File, output *os.File) {
	scanner := bufio.NewScanner(input)
	encoder := json.NewEncoder(output)

	ch := make(chan []TopSolution)
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
