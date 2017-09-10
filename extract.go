package main

import (
	"bufio"

	"github.com/pedmiston/foldit/solution"
)

// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent goroutines
var tokens = make(chan struct{}, 20)

func loadSolutions(in *bufio.Scanner) (out chan *solution.Solution, n int) {
	out = make(chan *solution.Solution)
	for in.Scan() {
		go func(path string) {
			tokens <- struct{}{}
			s, err := solution.New(path)
			if err != nil {
			  fmt.Fprintf(os.Stderr, "error loading solution: %s", path)
			}
			<-tokens
			out <- s
		}(in.Text())
		n++
	}
	return out, n
}

func writeSolutions(in chan *solution.Solution, n int) {
  for
}
