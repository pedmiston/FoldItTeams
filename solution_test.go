package main

import "testing"

func TestGettingPuzzleIDFromFilename(t *testing.T) {
	puzzleID, _ := readPuzzleIDFromFilename(fullSolution)
	if puzzleID != 2003996 {
		t.Error("Expected to extract puzzleID=2003996 but got", puzzleID)
	}
}
