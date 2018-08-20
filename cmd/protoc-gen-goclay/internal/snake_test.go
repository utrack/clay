package internal_test

import (
	"strings"
	"testing"

	"github.com/utrack/clay/v2/cmd/protoc-gen-goclay/internal"
)

type SnakeTest struct {
	input  string
	output string
}

var tests = []SnakeTest{
	{"a", "a"},
	{"snake", "snake"},
	{"A", "a"},
	{"ID", "id"},
	{"MOTD", "motd"},
	{"Snake", "snake"},
	{"SnakeID", "snake_id"},
	{"SnakeIDGoogle", "snake_id_google"},
	{"LinuxMOTD", "linux_motd"},
	{"OMGWTFBBQ", "omgwtfbbq"},
	{"omg_wtf_bbq", "omg_wtf_bbq"},
}

func TestFromString(t *testing.T) {
	for _, test := range tests {
		if internal.ToSnake(test.input) != test.output {
			t.Errorf(`ToSnake("%s"), wanted "%s", got "%s"`, test.input, test.output, internal.ToSnake(test.input))
		}
	}
}

var benchmarks = []string{
	"a",
	"snake",
	"A",
	"Snake",
	"SnakeTest",
	"SnakeID",
	"SnakeIDGoogle",
	"LinuxMOTD",
	"OMGWTFBBQ",
	"omg_wtf_bbq",
}

func BenchmarkToSnake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, input := range benchmarks {
			internal.ToSnake(input)
		}
	}
}

func BenchmarkToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, input := range benchmarks {
			strings.ToLower(input)
		}
	}
}
