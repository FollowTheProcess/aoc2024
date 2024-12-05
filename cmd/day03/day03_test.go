package main

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/test"
)

const testInput = "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"

func TestParseMul(t *testing.T) {
	tests := []struct {
		name    string // Name of the test case
		input   string // The input to parse
		want    Mul    // Expected mul instruction
		wantErr bool   // Do we want a parse error
	}{
		{
			name:    "empty",
			input:   "",
			want:    Mul{},
			wantErr: true,
		},
		{
			name:    "no muls",
			input:   "random words here but nothing we want",
			want:    Mul{},
			wantErr: true,
		},
		{
			name:  "valid",
			input: "mul(5,3)",
			want:  Mul{X: 5, Y: 3},
		},
		{
			name:  "3 digits",
			input: "mul(555,333)",
			want:  Mul{X: 555, Y: 333},
		},
		{
			name:    "invalid opening bracket",
			input:   "mul[5,3)",
			want:    Mul{},
			wantErr: true,
		},
		{
			name:    "invalid closing bracket",
			input:   "mul(5,3]",
			want:    Mul{},
			wantErr: true,
		},
		{
			name:    "no comma",
			input:   "mul(53)",
			want:    Mul{},
			wantErr: true,
		},
		{
			name:    "spaces",
			input:   "mul( 5 3 )",
			want:    Mul{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMul(tt.input)
			test.WantErr(t, err, tt.wantErr)
			test.Equal(t, got, tt.want)
		})
	}
}

func TestParseMuls(t *testing.T) {
	got, err := parseMuls(testInput)
	test.Ok(t, err)

	want := []Mul{
		{X: 2, Y: 4},
		{X: 5, Y: 5},
		{X: 11, Y: 8},
		{X: 8, Y: 5},
	}

	test.EqualFunc(t, got, want, slices.Equal)
}

func TestPart1Example(t *testing.T) {
	muls, err := parseMuls(testInput)
	test.Ok(t, err)

	sum := 0
	for _, mul := range muls {
		sum += mul.Do()
	}

	test.Equal(t, sum, 161) // Wrong answer for part 1 example
}
