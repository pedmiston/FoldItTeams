package irdata

import (
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
