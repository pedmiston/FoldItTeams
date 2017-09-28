package main

import (
	"bufio"
	"io"
)

// A Result contains the data extracted from a solution file and any error.
type Result struct {
	Data IRData
	Err  error
}

// tokens is a counting semaphore used to
// enforce a limited number of open files.
var tokens = make(chan struct{}, 20)

// ReadAll reads solution files concurrently.
func ReadAll(in io.Reader) (out chan Result, n int) {
	scanner := bufio.NewScanner(in)
	out = make(chan Result)
	for scanner.Scan() {
		go func(f string) {
			tokens <- struct{}{} // acquire a token
			data, err := Read(f)
			<-tokens // release the token
			result := Result{data, err}
			out <- result
		}(scanner.Text())
		n++
	}
	return out, n
}
