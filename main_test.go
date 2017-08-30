package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestAcceptsInputFromFile(t *testing.T) {
	tmpInputFile := createTmpFile()
	defer os.Remove(tmpInputFile.Name())
	tmpInputFile.WriteString(fullSolution)
	tmpInputFile.Close()
	tmpInputFile, _ = os.Open(tmpInputFile.Name())

	if inputLines := countLinesInFile(tmpInputFile.Name()); inputLines != 1 {
		t.Error("Problem writing input lines to file")
	}

	tmpOutputFile := createTmpFile()
	defer os.Remove(tmpOutputFile.Name())

	writeTopSolutionsToJSON(tmpInputFile, tmpOutputFile)

	if outputLines := countLinesInFile(tmpOutputFile.Name()); outputLines != 1 {
		t.Error("Expected 1 output lines but got", outputLines)
	}

}

func TestSkippingBadFiles(t *testing.T) {
	tmpInputFile := createTmpFile()
	defer os.Remove(tmpInputFile.Name())
	tmpInputFile.WriteString(fullSolution + "\n" + badFilename)
	tmpInputFile.Close()
	tmpInputFile, _ = os.Open(tmpInputFile.Name())

	if inputLines := countLinesInFile(tmpInputFile.Name()); inputLines != 2 {
		t.Error("Problem writing input lines to file")
	}

	tmpOutputFile := createTmpFile()
	defer os.Remove(tmpOutputFile.Name())

	writeTopSolutionsToJSON(tmpInputFile, tmpOutputFile)

	if outputLines := countLinesInFile(tmpOutputFile.Name()); outputLines != 2 {
		t.Error("Expected 2 output lines but got", outputLines)
	}
}

func createTmpFile() *os.File {
	tmp, err := ioutil.TempFile("", "foldit")
	if err != nil {
		log.Fatal(err)
	}
	return tmp
}

func countLinesInFile(fname string) int {
	f, _ := os.Open(fname)
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	counter := 0
	for scanner.Scan() {
		counter++
	}
	return counter
}
