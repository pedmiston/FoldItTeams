package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pedmiston/foldit/solution"
)

type result struct {
	s   *solution.Solution
	err error
}

func loadSolutions(in *bufio.Scanner) (out chan result, n int) {
	out = make(chan result)
	for in.Scan() {
		go readSolution(in.Text(), out)
		n++
	}
	return out, n
}

// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent goroutines
var tokens = make(chan struct{}, 20)

func readSolution(path string, dst chan result) {
	tokens <- struct{}{} // obtain token
	s, err := solution.New(path)
	r := result{s, err}
	dst <- r
	<-tokens // release token
}

func writeSolutions(results chan result, n int, encoder *json.Encoder) {
	for i := 0; i < n; i++ {
		r := <-results
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "%s,%v\n", r.s.Filepath, r.err)
			continue
		}
		if err := encoder.Encode(r.s); err != nil {
			fmt.Fprintf(os.Stderr, "%s,%v\n", r.s.Filepath, err)
		}
	}
}
