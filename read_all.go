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

// ReadAll reads solution files concurrently.
// TODO: Place bounds on concurrency
func ReadAll(in io.Reader) (out chan Result, n int) {
	scanner := bufio.NewScanner(in)
	out = make(chan Result)
	for scanner.Scan() {
		go func(f string) {
			data, err := Read(f)
			result := Result{data, err}
			out <- result
		}(scanner.Text())
		n++
	}
	return out, n
}
