package main

// A TopSolution is a solution with ranking information
type TopSolution struct {
	Solution
	RankType string
	Rank     int
}
