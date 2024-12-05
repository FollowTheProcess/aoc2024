/*
--- Day 2: Red-Nosed Reports ---

Fortunately, the first location The Historians want to search isn't a long walk from the Chief Historian's office.

While the Red-Nosed Reindeer nuclear fusion/fission plant appears to contain no sign of the Chief Historian, the engineers there run up to you as soon
as they see you. Apparently, they still talk about the time Rudolph was saved through molecular synthesis from a single electron.

They're quick to add that - since you're already here - they'd really appreciate your help analyzing some unusual data from the Red-Nosed reactor.
You turn to check if The Historians are waiting for you, but they seem to have already divided into groups that are currently searching every corner
of the facility. You offer to help with the unusual data.

The unusual data (your puzzle input) consists of many reports, one report per line. Each report is a list of numbers called levels that are separated by spaces.
For example:

7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9

This example data contains six reports each containing five levels.

The engineers are trying to figure out which reports are safe.

The Red-Nosed reactor safety systems can only tolerate levels that are either gradually increasing or gradually decreasing. So, a report only counts as safe if both of the following are true:

The levels are either all increasing or all decreasing.
Any two adjacent levels differ by at least one and at most three.
In the example above, the reports can be found safe or unsafe by checking those rules:

7 6 4 2 1: Safe because the levels are all decreasing by 1 or 2.
1 2 7 8 9: Unsafe because 2 7 is an increase of 5.
9 7 6 2 1: Unsafe because 6 2 is a decrease of 4.
1 3 2 4 5: Unsafe because 1 3 is increasing but 3 2 is decreasing.
8 6 4 4 1: Unsafe because 4 4 is neither an increase or a decrease.
1 3 6 7 9: Safe because the levels are all increasing by 1, 2, or 3.

So, in this example, 2 reports are safe.

Analyze the unusual data from the engineers. How many reports are safe?

--- Part Two ---

The engineers are surprised by the low number of safe reports until they realize they forgot to tell you about the Problem Dampener.

The Problem Dampener is a reactor-mounted module that lets the reactor safety systems tolerate a single bad level in what would otherwise be a safe report.
It's like the bad level never happened!

Now, the same rules apply as before, except if removing a single level from an unsafe report would make it safe, the report instead counts as safe.

More of the above example's reports are now safe:

7 6 4 2 1: Safe without removing any level.
1 2 7 8 9: Unsafe regardless of which level is removed.
9 7 6 2 1: Unsafe regardless of which level is removed.
1 3 2 4 5: Safe by removing the second level, 3.
8 6 4 4 1: Safe by removing the third level, 4.
1 3 6 7 9: Safe without removing any level.

Thanks to the Problem Dampener, 4 reports are actually safe!

Update your analysis by handling situations where the Problem Dampener can remove a single level from unsafe reports. How many reports are now safe?
*/

package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/FollowTheProcess/msg"
)

//go:embed day02.txt
var input string

func main() {
	if err := run(input); err != nil {
		msg.Error("%v", err)
		os.Exit(1)
	}
}

func run(input string) error {
	reports, err := parseInput(input)
	if err != nil {
		return err
	}

	fmt.Printf("Part 1: %d\n", countSafe(reports))
	fmt.Printf("Part 2: %d\n", countSafeRelaxed(reports))

	return nil
}

// parseInput parses a list of Reports from the puzzle input.
func parseInput(input string) ([]Report, error) {
	input = strings.TrimSpace(input)

	scanner := bufio.NewScanner(strings.NewReader(input))
	lineNo := 1

	var reports []Report

	for scanner.Scan() {
		line := scanner.Text()
		rawLevels := strings.Split(line, " ")

		report := make([]int, 0, len(rawLevels))

		for _, rawLevel := range rawLevels {
			level, err := strconv.Atoi(rawLevel)
			if err != nil {
				return nil, fmt.Errorf("bad level (%q in %q) on line %d: %w", rawLevel, line, lineNo, err)
			}
			report = append(report, level)
		}

		reports = append(reports, report)

		lineNo++
	}

	return reports, nil
}

// countSafe returns the number of reports that are safe.
func countSafe(reports []Report) int {
	safe := 0
	for _, report := range reports {
		if report.IsSafe() {
			safe++
		}
	}

	return safe
}

// countSafeRelaxed returns the number of reports that are safe with
// the problem dampener taken into account.
func countSafeRelaxed(reports []Report) int {
	safe := 0
	for _, report := range reports {
		if report.IsSafeRelaxed() {
			safe++
		}
	}

	return safe
}

// Report repesents a report from the red-nosed reactor.
type Report []int

// IsSafe reports whether the report is safe according to the puzzle criteria.
//
// A report only counts as safe if both of the following are true:
//   - The levels are either all increasing or all decreasing.
//   - Any two adjacent levels differ by at least one and at most three.
func (r Report) IsSafe() bool {
	// Yes this iterates through the slice 3 times but it's advent of code
	// so who cares
	if !r.allIncreasing() && !r.allDecreasing() {
		return false
	}

	if !r.differenceAllowed() {
		return false
	}

	return true
}

// IsSafeRelaxed is like IsSafe but takes the problem dampener into account.
func (r Report) IsSafeRelaxed() bool {
	if r.IsSafe() {
		// If it's safe already it's obviously still safe
		return true
	}

	// Remove one at a time and test for safety
	for i := 0; i < len(r); i++ {
		removed := append(append([]int{}, r[0:i]...), r[i+1:]...)
		report := Report(removed)
		if report.IsSafe() {
			return true
		}
	}

	// Still not safe
	return false
}

// allDecreasing reports whether the Report contains values that are
// always decreasing e.g. 5, 4, 3, 2, 1.
func (r Report) allDecreasing() bool {
	for i := 1; i < len(r); i++ {
		if r[i] >= r[i-1] {
			// Next element is greater than or equal to the previous
			// so it can't be all decreasing
			return false
		}
	}

	return true
}

// allIncreasing reports whether the Report contains values that are
// always increasing e.g. 1, 2, 3, 4, 5.
func (r Report) allIncreasing() bool {
	for i := 1; i < len(r); i++ {
		if r[i] <= r[i-1] {
			// Next element is less than or equal to the previous
			// so it can't be all increasing
			return false
		}
	}

	return true
}

// differenceAllowed reports whether the difference between any two adjacent levels
// is allowed.
func (r Report) differenceAllowed() bool {
	for i := 1; i < len(r); i++ {
		current := r[i]
		previous := r[i-1]

		diff := int(math.Abs(float64(current - previous)))

		if diff > 3 || diff < 1 {
			return false
		}
	}

	return true
}
