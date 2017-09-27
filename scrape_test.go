package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func BenchmarkScrape(b *testing.B) {
	tmpDir, paths := replicate("testdata/small_solution.pdb", 10)
	defer os.RemoveAll(tmpDir)

	var dst, errDst io.Writer
	dst = new(bytes.Buffer)
	errDst = new(bytes.Buffer)

	b.Run("N=10", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Scrape(paths, dst, errDst)
		}
	})
}
