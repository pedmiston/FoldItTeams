package main

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	wd, _        = os.Getwd()
	solutionDir  = "testdata"
	topSolution  = "solution_2003996/top/solution_bid_0004_0000288912_0002003867_0372197993.ir_solution.pdb"
	fullSolution = filepath.Join(wd, solutionDir, topSolution)
	badFilename  = "solution_1/bottom/a_file.pdb"
)

func TestGettingRankFromFilename(t *testing.T) {
	_, rank, _ := getRankFromFilename(fullSolution)
	if rank != 4 {
		t.Error("Expected to extract rank=4 but got", rank)
	}
}
