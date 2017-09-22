package irdata

import (
	"bufio"
	"io"
)

// A Result comprises the data extracted from a solution file and any error.
type Result struct {
	Data IRDATA
	Err  error
}

// Load reads solution files concurrently.
func Load(in io.Reader) (out chan Result, n int) {
	scanner := bufio.NewScanner(in)
	out = make(chan Result)
	for scanner.Scan() {
		go func(pdb string) {
			data, err := Read(pdb)
			result := Result{data, err}
			out <- result
		}(scanner.Text())
		n++
	}
	return out, n
}
