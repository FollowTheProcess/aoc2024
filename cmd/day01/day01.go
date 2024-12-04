package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/FollowTheProcess/collections/counter"
	"github.com/FollowTheProcess/msg"
)

//go:embed day01.txt
var input string

func main() {
	if err := run(input); err != nil {
		msg.Error("%v", err)
		os.Exit(1)
	}
}

func run(input string) error {
	left, right, err := parseInput(input)
	if err != nil {
		return err
	}

	distance := totalDistance(left, right)
	similarity := similarityScore(left, right)

	fmt.Printf("Part 1: %d\n", distance)
	fmt.Printf("Part2: %d\n", similarity)

	return nil
}

// totalDistance calculates the total distance between the two lists, where
// the distance is the sum of the differences between each element in the sorted
// lists.
func totalDistance(left, right []int) int {
	// Sort each so we can pair elements up in order
	slices.Sort(left)
	slices.Sort(right)

	// The sum of all the diffs
	sum := 0
	for index, leftValue := range left {
		rightValue := right[index]
		diff := math.Abs(float64(leftValue - rightValue))
		sum += int(diff)
	}

	return sum
}

// similarityScore calculates the similarity score of the two lists by iterating through
// the numbers in the left list, multiplying them by the number of times they occur in the
// right list, and summing all this together.
func similarityScore(left, right []int) int {
	// Counts of everything in the right hand list
	counts := counter.From(right)

	sum := 0
	for _, leftValue := range left {
		// Multiply the number by the number of times it occurs
		// in the right list
		similarity := leftValue * counts.Count(leftValue)
		sum += similarity
	}

	return sum
}

// parseInput parses the raw input text into two lists of integers representing
// the left list and the right list.
func parseInput(input string) (left, right []int, err error) {
	input = strings.TrimSpace(input)
	scanner := bufio.NewScanner(strings.NewReader(input))

	lineNo := 1

	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		leftRaw, rightRaw, ok := strings.Cut(line, "   ")
		if !ok {
			return nil, nil, fmt.Errorf("Bad line %s", line)
		}

		leftInt, err := strconv.Atoi(leftRaw)
		if err != nil {
			return nil, nil, fmt.Errorf("Bad numeric value in left list %s on line %d: %w", leftRaw, lineNo, err)
		}

		rightInt, err := strconv.Atoi(rightRaw)
		if err != nil {
			return nil, nil, fmt.Errorf("Bad numeric value in right list %s on line %d: %w", rightRaw, lineNo, err)
		}

		left = append(left, leftInt)
		right = append(right, rightInt)
	}

	return left, right, nil
}
