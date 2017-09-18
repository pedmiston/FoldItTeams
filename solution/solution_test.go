package solution

import "testing"

var testSolution = "testdata/test_solution_1.pdb"

func TestNewSolution(t *testing.T) {
	s, err := New(testSolution)
	if err != nil {
		t.Error("Creating new solution", err)
	}
	expectedSolutionID := 372197993
	if s.SolutionID != expectedSolutionID {
		t.Errorf("expected s.SolutionID == %v, got %v", expectedSolutionID, s.SolutionID)
	}
}
