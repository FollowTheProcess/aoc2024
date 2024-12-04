package main

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/test"
)

const testInput = `
3   4
4   3
2   5
1   3
3   9
3   3
`

func TestParseInput(t *testing.T) {
	left, right, err := parseInput(testInput)
	test.Ok(t, err)

	wantLeft := []int{3, 4, 2, 1, 3, 3}
	wantRight := []int{4, 3, 5, 3, 9, 3}

	test.EqualFunc(t, left, wantLeft, slices.Equal)
	test.EqualFunc(t, right, wantRight, slices.Equal)
}

func TestTotalDifference(t *testing.T) {
	left := []int{3, 4, 2, 1, 3, 3}
	right := []int{4, 3, 5, 3, 9, 3}

	want := 11

	test.Equal(t, totalDistance(left, right), want)
}

func TestSimilarityScore(t *testing.T) {
	left := []int{3, 4, 2, 1, 3, 3}
	right := []int{4, 3, 5, 3, 9, 3}

	want := 31

	test.Equal(t, similarityScore(left, right), want)
}
