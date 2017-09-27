package main

import "testing"

func TestSuccessfulRead(t *testing.T) {
	_, err := Read("testdata/small_solution.pdb")
	if err != nil {
		t.Errorf("error reading 'testdata/small_solution.pdb': %v", err)
	}
}

func TestReadValues(t *testing.T) {
	data, _ := Read("testdata/small_solution.pdb")

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
		gotString := got.(string)
		if !ok {
			t.Errorf("IRDATA field '%v' not extracted", test.key)
		} else {
			if gotString != test.expected {
				t.Errorf("expected %v = %v, got %v", test.key, test.expected, gotString)
			}
		}
	}
}

func TestReadStoresFilepath(t *testing.T) {
	src := "testdata/empty_solution.pdb"
	data, _ := Read(src)
	if data.Filepath() != src {
		t.Errorf("expected Read to store filepath as attribute")
	}
}

func TestReadMultipleUsers(t *testing.T) {
	data, _ := Read("testdata/multiple_users.pdb")
	pdls, ok := data["PDL"]
	if !ok {
		t.Fatal("IRDATA field 'PDL' not found")
	}
	slice := pdls.([]string)
	if len(slice) != 3 {
		t.Fatalf("expected to pull 3 PDLs from 'testdata/multiple_users.pdb' but got %v", len(slice))
	}
}

func BenchmarkRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Read("testdata/real_solution.pdb")
	}
}
