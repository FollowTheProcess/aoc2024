/*
--- Day 3: Mull It Over ---

"Our computers are having issues, so I have no idea if we have any Chief Historians in stock! You're welcome to check the warehouse, though,"
says the mildly flustered shopkeeper at the North Pole Toboggan Rental Shop. The Historians head out to take a look.

The shopkeeper turns to you. "Any chance you can see why our computers are having issues again?"

The computer appears to be trying to run a program, but its memory (your puzzle input) is corrupted. All of the instructions have been jumbled up!

It seems like the goal of the program is just to multiply some numbers. It does that with instructions like mul(X,Y), where X and Y are each 1-3 digit numbers.
For instance, mul(44,46) multiplies 44 by 46 to get a result of 2024. Similarly, mul(123,4) would multiply 123 by 4.

However, because the program's memory has been corrupted, there are also many invalid characters that should be ignored, even if they look like
part of a mul instruction. Sequences like mul(4*, mul(6,9!, ?(12,34), or mul ( 2 , 4 ) do nothing.

For example, consider the following section of corrupted memory:

xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))

Only the four highlighted sections are real mul instructions. Adding up the result of each instruction produces 161 (2*4 + 5*5 + 11*8 + 8*5).

Scan the corrupted memory for uncorrupted mul instructions. What do you get if you add up all of the results of the multiplications?

--- Part Two ---

As you scan through the corrupted memory, you notice that some of the conditional statements are also still intact.
If you handle some of the uncorrupted conditional statements in the program, you might be able to get an even more accurate result.

There are two new instructions you'll need to handle:

The do() instruction enables future mul instructions.
The don't() instruction disables future mul instructions.

Only the most recent do() or don't() instruction applies. At the beginning of the program, mul instructions are enabled.

For example:

xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))

This corrupted memory is similar to the example from before, but this time the mul(5,5) and mul(11,8) instructions are disabled because there
is a don't() instruction before them. The other mul instructions function normally, including the one at the end that gets re-enabled by a do() instruction.

This time, the sum of the results is 48 (2*4 + 8*5).

Handle the new instructions; what do you get if you add up all of the results of just the enabled multiplications?
*/

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/parser"
)

const (
	mulRegexRaw = `mul\((\d{1,3}),(\d{1,3})\)`
	allRegexRaw = `mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)` // Find any
)

//go:embed day03.txt
var input string

var (
	mulRegex = regexp.MustCompile(mulRegexRaw)
	allRegex = regexp.MustCompile(allRegexRaw)
)

func main() {
	if err := run(); err != nil {
		msg.Error("%v", err)
		os.Exit(1)
	}
}

func run() error {
	muls, err := parseMuls(input)
	if err != nil {
		return err
	}

	sum := 0
	for _, mul := range muls {
		sum += mul.Do()
	}

	fmt.Printf("Part 1: %d\n", sum)

	enabledMuls, err := parseEnabledMuls(input)
	if err != nil {
		return err
	}

	enabledSum := 0
	for _, enabled := range enabledMuls {
		enabledSum += enabled.Do()
	}

	fmt.Printf("Part 2: %d\n", enabledSum)
	return nil
}

// Mul represents a multiply instruction.
type Mul struct {
	X     int // The left operand
	Y     int // The right operand
	Start int // The start position of this mul in the input string
}

// Do performs the multiplication, returning the answer.
func (m Mul) Do() int {
	return m.X * m.Y
}

// parseMuls parses 1 or more mul instructions from the input string.
func parseMuls(input string) ([]Mul, error) {
	spans := mulRegex.FindAllStringIndex(input, -1)
	if len(spans) == 0 {
		return nil, errors.New("no muls found")
	}

	parsed := make([]Mul, 0, len(spans))
	for _, span := range spans {
		start := span[0]
		end := span[1]
		m, err := parseMul(input[start:end])
		if err != nil {
			return nil, err
		}
		m.Start = start
		parsed = append(parsed, m)
	}

	return parsed, nil
}

// parseEnabledMuls parses 1 or more mul instructions from the input string but only
// if they fall in between a do and a don't statement so are enabled.
func parseEnabledMuls(input string) ([]Mul, error) {
	matches := allRegex.FindAllString(input, -1)
	if len(matches) == 0 {
		return nil, errors.New("no mul, do or don't found")
	}

	var parsed []Mul

	// Initial state is enabled
	enabled := true

	for _, match := range matches {
		switch match {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			// Must be a mul, only add it if we're currently enabled
			if enabled {
				m, err := parseMul(match)
				if err != nil {
					return nil, err
				}
				parsed = append(parsed, m)
			}
		}
	}

	return parsed, nil
}

// parseMul parses a single mul instruction string.
func parseMul(raw string) (Mul, error) {
	// We don't actually care about the mul keyword, as long
	// as it's there
	keyword, rest, err := parser.Exact("mul")(raw)
	if err != nil {
		return Mul{}, err
	}

	if keyword != "mul" {
		return Mul{}, fmt.Errorf("expected keyword 'mul' got %q", keyword)
	}

	// Next up should be a left bracket '(', again we don't capture it
	// just need to know it's there
	leftBracket, rest, err := parser.Exact("(")(rest)
	if err != nil {
		return Mul{}, err
	}

	if leftBracket != "(" {
		return Mul{}, fmt.Errorf("expected '(', got %q", leftBracket)
	}

	// Now find the operands
	leftOperand, rest, err := parser.TakeWhile(unicode.IsDigit)(rest)
	if err != nil {
		return Mul{}, err
	}

	if len(leftOperand) == 0 {
		return Mul{}, errors.New("expected digit for left operand")
	}

	// Should be a comma now
	comma, rest, err := parser.Exact(",")(rest)
	if err != nil {
		return Mul{}, err
	}

	if comma != "," {
		return Mul{}, fmt.Errorf("expected two numeric operands separated by a comma, got %q", comma)
	}

	rightOperand, rest, err := parser.TakeWhile(unicode.IsDigit)(rest)
	if err != nil {
		return Mul{}, err
	}

	if len(rightOperand) == 0 {
		return Mul{}, errors.New("expected digit for right operand")
	}

	// Now just the closing bracket and we're done
	rightBracket, _, err := parser.Exact(")")(rest)
	if err != nil {
		return Mul{}, err
	}

	if rightBracket != ")" {
		return Mul{}, fmt.Errorf("expected ')', got %q", rightBracket)
	}

	// We have everything we need, just type conversions
	left, err := strconv.Atoi(leftOperand)
	if err != nil {
		return Mul{}, fmt.Errorf("bad left operand %q: %w", leftOperand, err)
	}

	right, err := strconv.Atoi(rightOperand)
	if err != nil {
		return Mul{}, fmt.Errorf("bad right operand %q: %w", rightOperand, err)
	}

	return Mul{X: left, Y: right}, nil
}
