package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	data, err := Read("testdata/small_solution.pdb")
	if err != nil {
		t.Errorf("error reading 'testdata/small_solution.pdb': %v", err)
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"TITLE", "Status Solution"},
		{"PID", "2002990"},
		{"DESCRIPTION", "Generated on Mon Oct 24 19:01:45 2016."},
	}

	for _, test := range tests {
		got, ok := data[test.key]
		if !ok {
			t.Errorf("IRDATA field '%v' not extracted", test.key)
		} else {
			if len(got) == 1 {
				if g := got[0]; g != test.expected {
					t.Errorf("expected %v = %v, got %v", test.key, test.expected, g)
				}
			} else {
				t.Errorf("expected a single value, got: %v", got)
			}
		}
	}
}

func TestReadStoresFilepath(t *testing.T) {
	src := "testdata/small_solution.pdb"
	data, _ := Read(src)
	v, ok := data["FILEPATH"]
	if !ok || v[0] != src {
		t.Errorf("expected Read to store filepath as attribute")
	}
}

func TestReadMultipleUsers(t *testing.T) {
	data, err := Read("testdata/multiple_users.pdb")
	if err != nil {
		t.Errorf("error reading 'testdata/multiple_users.pdb': %v", err)
	}
	pdls, ok := data["PDL"]
	if !ok {
		t.Fatal("IRDATA field 'PDL' not found")
	}
	if len(pdls) != 3 {
		t.Fatalf("expected to pull 3 PDLs from 'testdata/multiple_users.pdb' but got %v", len(pdls))
	}
}

func BenchmarkRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Read("testdata/real_solution.pdb")
	}
}

func TestLoad(t *testing.T) {
	in := strings.NewReader("testdata/small_solution.pdb\n")
	ch, n := Load(in)
	if n != 1 {
		t.Errorf("Expected to load 1 solution, instead loaded %v", n)
	}
	r := <-ch
	if r.Err != nil {
		t.Errorf("Expected to load 1 solution without error, got %v", r.Err)
	}
	title, ok := r.Data["TITLE"]
	if !ok || title[0] != "Status Solution" {
		t.Errorf("Expected to load Title = Status Solution, got %v", title[0])
	}
}

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

func BenchmarkWrite(b *testing.B) {
	tmpDir, paths := replicate("testdata/small_solution.pdb", 10)
	defer os.RemoveAll(tmpDir)

	var dst io.Writer
	dst = new(bytes.Buffer)

	b.Run("Write=10", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Write(paths, dst)
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
