package main

import (
	"bufio"
	"os"
)

// Read returns a map of all IRData fields in a solution file.
func Read(f string) (IRData, error) {
	data := make(IRData)
	data["FILEPATH"] = f

	in, err := os.Open(f)
	if err != nil {
		return data, err
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if reIRLine.MatchString(line) {
			err = data.Append(line)
			if err != nil {
				return data, err
			}
		}
	}

	return data, nil
}
