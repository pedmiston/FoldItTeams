package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

var reIRDATA = regexp.MustCompile(`^IRDATA ([^ ]*) (.*)$`)

// A IRDATA maps data fields to one or more lines of a solution file.
type IRDATA map[string][]string

// A Result comprises the data extracted from a solution file and any error.
type Result struct {
	Data IRDATA
	Err  error
}

// Read returns a map of all IRDATA fields in a solution file.
func Read(pdb string) (IRDATA, error) {
	var in *os.File
	var err error

	in, err = os.Open(pdb)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	data := make(IRDATA)
	data["FILEPATH"] = []string{pdb}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if reIRDATA.MatchString(line) {
			matches := reIRDATA.FindStringSubmatch(line)
			if len(matches) != 3 {
				fmt.Fprintf(os.Stderr, "%s,%s", pdb, line)
				continue
			}
			key, line := matches[1], matches[2]

			prev, ok := data[key]
			if !ok {
				data[key] = []string{line}
			} else {
				data[key] = append(prev, line)
			}
		}
	}

	return data, nil
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

// Write reads solution files, selects all IRDATA, and writes to JSON.
func Write(in io.Reader, dst io.Writer) {
	// Load a channel of IRDATA from the input scanner.
	out, n := Load(in)

	// Pull data from the channel and encode it to JSON.
	encoder := json.NewEncoder(dst)
	for i := 0; i < n; i++ {
		r := <-out
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "%s,%v\n", r.Data["FILEPATH"], r.Err)
			continue
		}
		if err := encoder.Encode(r.Data); err != nil {
			fmt.Fprintf(os.Stderr, "%s,%v\n", r.Data["FILEPATH"], r.Err)
		}
	}
}
