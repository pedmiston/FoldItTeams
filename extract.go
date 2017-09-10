package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pedmiston/foldit/solution"
)

func loadSolutions(in *bufio.Scanner) (out chan *solution.Solution, n int) {
	out = make(chan *solution.Solution)
	for in.Scan() {
		go readSolution(in.Text(), out)
		n++
	}
	return out, n
}

// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent goroutines
var tokens = make(chan struct{}, 20)

func readSolution(path string, dst chan *solution.Solution) {
	tokens <- struct{}{} // obtain token
	s, err := solution.New(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading solution %s: %v", path, err)
	}
	dst <- s
	<-tokens // release token
}
