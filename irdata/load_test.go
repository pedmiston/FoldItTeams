package irdata

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func BenchmarkLoad(b *testing.B) {
	tmpDir, paths := replicate("testdata/small_solution.pdb", 10)
	defer os.RemoveAll(tmpDir)
	b.Run("Load=10", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			ch, n := Load(paths)
			for i := 0; i < n; i++ {
				r := <-ch
				if r.Err != nil {
					b.Log(r.Err)
				}
			}
		}
	})
}

func replicate(src string, n int) (string, io.Reader) {
	// Find out how big the src is
	info, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	// Create a byte slice big enough to hold the solution
	solution := make([]byte, info.Size())

	// Open the solution and read it into the byte slice
	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Read(solution)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// Create the directory for solution files
	tmpDir, err := ioutil.TempDir(".", strings.Split(src, ".")[0]+"_")
	if err != nil {
		log.Fatal(err)
	}

	// Create new files for each replicate solution,
	// write the solution byte slice to the file,
	// and append the filename to a string of paths.
	var paths string
	for i := 0; i < n; i++ {
		dst, err := ioutil.TempFile(tmpDir, "")
		if err != nil {
			log.Fatal(err)
		}

		_, err = dst.Write(solution)
		if err != nil {
			log.Fatal(err)
		}

		dst.Close()
		paths += dst.Name() + "\n"
	}

	return tmpDir, strings.NewReader(paths)
}
