package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestExtract(t *testing.T) {
	replicates := 100
	pathsFilepath, err := replicateSolutions("solution/testdata/test_solution_1.pdb", "testdata", replicates)
	if err != nil {
		t.Errorf("error replicating solutions: %v", err)
	}
	defer os.RemoveAll("testdata")

	pathsFile, err := os.Open(pathsFilepath)
	if err != nil {
		t.Errorf("couldn't open paths file")
	}
	defer pathsFile.Close()

	scanner := bufio.NewScanner(pathsFile)
	ch, n := loadSolutions(scanner)

	r := <-ch

	for i := 1; i < n; i++ {
		<-ch
	}

	if r.s.PuzzleID != 2003996 {
		t.Error("solution not parsed correctly")
	}

	if replicates != n {
		t.Errorf("incorrect number of solutions, expected %v, got %v", replicates, n)
	}
}

func replicateSolutions(srcPath, dstDir string, n int) (string, error) {
	// Find out how big the solution is
	info, err := os.Stat(srcPath)
	if err != nil {
		return "", err
	}

	// Create a byte slice big enough to hold the solution
	solution := make([]byte, info.Size())

	// Open the solution and read it into the byte slice
	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	_, err = src.Read(solution)
	if err != nil {
		return "", err
	}
	src.Close()

	// Create the directory for solution files
	err = os.Mkdir(dstDir, 0777)
	if err != nil {
		return "", err
	}

	// Create the file to hold solution paths
	pathsFilepath := path.Join(dstDir, "filepaths.txt")
	pathsFile, err := os.Create(pathsFilepath)
	if err != nil {
		return "", err
	}

	// Create new files for each replicate solution,
	// write the solution byte slice to the file,
	// and write the new file path to the paths file.
	for i := 0; i < n; i++ {
		newSolutionPath := path.Join(dstDir, fmt.Sprintf("solution_%v.pdb", i))

		dst, err := os.Create(newSolutionPath)
		if err != nil {
			return "", err
		}
		_, err = dst.Write(solution)
		if err != nil {
			return "", err
		}
		dst.Close()

		pathsFile.WriteString(newSolutionPath + "\n")
	}
	pathsFile.Close()

	return pathsFilepath, nil
}
