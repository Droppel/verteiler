package main

import (
	"testing"
)

func TestSolveSubSum(t *testing.T) {
	sum := 72
	set := []int{2, 1, 3, 5, 67}

	solved, subset := solveSubsetSum(Solution{}, sum, set)
	t.Log(solved)
	t.Log(subset)
	t.Fail()
}
