package main

import (
	"testing"

	"github.com/FollowTheProcess/test"
)

const testInput = `
7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9
`

func TestParse(t *testing.T) {
	reports, err := parseInput(testInput)
	test.Ok(t, err)

	want := []Report{
		{7, 6, 4, 2, 1},
		{1, 2, 7, 8, 9},
		{9, 7, 6, 2, 1},
		{1, 3, 2, 4, 5},
		{8, 6, 4, 4, 1},
		{1, 3, 6, 7, 9},
	}

	test.Diff(t, reports, want)
}

func TestCountSafe(t *testing.T) {
	reports := []Report{
		{7, 6, 4, 2, 1},
		{1, 2, 7, 8, 9},
		{9, 7, 6, 2, 1},
		{1, 3, 2, 4, 5},
		{8, 6, 4, 4, 1},
		{1, 3, 6, 7, 9},
	}

	want := 2

	got := countSafe(reports)

	test.Equal(t, got, want) // countSafe returned the wrong answer
}

func TestReport(t *testing.T) {
	tests := []struct {
		name          string // Name of the test case
		report        Report // The report under test
		isSafe        bool   // Whether the report should be considered safe
		allDecreasing bool   // Whether the report is all decreasing
		allIncreasing bool   // Whether the report is all increasing
	}{
		{
			name:          "example 1",
			report:        []int{7, 6, 4, 2, 1},
			isSafe:        true, // because the levels are all decreasing by 1 or 2
			allDecreasing: true,
		},
		{
			name:          "example 2",
			report:        []int{1, 2, 7, 8, 9},
			isSafe:        false, // because 2 7 is an increase of 5
			allIncreasing: true,
		},
		{
			name:          "example 3",
			report:        []int{9, 7, 6, 2, 1},
			isSafe:        false, // because 6 2 is a decrease of 4.
			allDecreasing: true,
		},
		{
			name:   "example 4",
			report: []int{1, 3, 2, 4, 5},
			isSafe: false, // because 1 3 is increasing but 3 2 is decreasing.
		},
		{
			name:          "example 5",
			report:        []int{8, 6, 4, 4, 1},
			isSafe:        false, // because 4 4 is neither an increase or a decrease.
			allDecreasing: false,
		},
		{
			name:          "example 6",
			report:        []int{1, 3, 6, 7, 9},
			isSafe:        true, // because the levels are all increasing by 1, 2, or 3.
			allIncreasing: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.report.allIncreasing(), tt.allIncreasing) // allIncreasing mismatch
			test.Equal(t, tt.report.allDecreasing(), tt.allDecreasing) // allDecreasing mismatch
			test.Equal(t, tt.report.IsSafe(), tt.isSafe)               // IsSafe() mismatch
		})
	}
}
