package main

import (
	"fmt"
	"regexp"
)

// An IRData maps data field keys to one or more lines of a solution file.
type IRData map[string]interface{}

// reIRLine matches an IRDATA line with groups for key and value
var reIRLine = regexp.MustCompile(`^IRDATA ([^ ]*) (.*)$`)

// Append adds new data from an IRDATA line in the solution file.
func (data IRData) Append(line string) error {
	matches := reIRLine.FindStringSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("expected 3 submatch groups")
	}

	key, line := matches[1], matches[2]

	prev, ok := data[key]
	if !ok {
		data[key] = line
	} else {
		switch prevTyped := prev.(type) {
		case string:
			data[key] = []string{prevTyped, line}
		case []string:
			data[key] = append(prevTyped, line)
		default:
			return fmt.Errorf("unknown dtype %v", prevTyped)
		}
	}

	return nil
}

// Filepath returns the solution filepath
func (data IRData) Filepath() string {
	f, _ := data["FILEPATH"]
	return f.(string)
}
