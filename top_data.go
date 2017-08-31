package main

import "strconv"

type TopData struct {
	Scores  [][]string
	Actions [][]string
	History [][]string
}

func NewTopData(topSolution *TopSolution) *TopData {
	topData := &TopData{
		Scores:  topSolution.getScores(),
		Actions: topSolution.getActions(),
		History: topSolution.getHistory(),
	}

	return topData
}

func (t *TopSolution) getScores() [][]string {
	return [][]string{[]string{
		strconv.Itoa(t.PuzzleID),
		strconv.Itoa(t.UserID),
		strconv.Itoa(t.GroupID),
		strconv.FormatFloat(t.Score, 'f', -1, 64),
		t.Filename,
	}}
}

func (t *TopSolution) getActions() [][]string {
	var row int
	dataframe := make([][]string, len(t.Actions))
	for action, count := range t.Actions {
		dataframe[row] = []string{
			strconv.Itoa(t.PuzzleID),
			strconv.Itoa(t.UserID),
			strconv.Itoa(t.GroupID),
			strconv.FormatFloat(t.Score, 'f', -1, 64),
			t.Filename,
			action,
			strconv.Itoa(count),
		}
		row++
	}
	return dataframe
}

func (t *TopSolution) getHistory() [][]string {
	var row int
	dataframe := make([][]string, len(t.History))
	for ix, id := range t.Actions {
		dataframe[row] = []string{
			strconv.Itoa(t.PuzzleID),
			strconv.Itoa(t.UserID),
			strconv.Itoa(t.GroupID),
			strconv.FormatFloat(t.Score, 'f', -1, 64),
			t.Filename,
			ix,
			strconv.Itoa(id),
		}
		row++
	}
	return dataframe
}
