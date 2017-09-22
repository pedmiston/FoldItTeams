package irdata

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var reIRDATA = regexp.MustCompile(`^IRDATA ([^ ]*) (.*)$`)

// IRDATA maps data fields to one or more lines in a solution file.
type IRDATA map[string][]string

// Read returns a map of all IRDATA fields in a solution file.
func Read(pdb string) (IRDATA, error) {
	var in *os.File
	var err error

	in, err = os.Open(pdb)
	if err != nil {
		return nil, fmt.Errorf("unable to open solution file: %v", err)
	}
	defer in.Close()

	data := make(map[string][]string)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if reIRDATA.MatchString(line) {
			matches := reIRDATA.FindStringSubmatch(line)
			if len(matches) != 3 {
				fmt.Fprintf(os.Stderr,
					"error: irdata line doesn't contain expected data:\n\t%s",
					line)
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
