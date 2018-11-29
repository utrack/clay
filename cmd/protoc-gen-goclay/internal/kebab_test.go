package internal

import (
	"testing"
)

type KebabTest struct {
	input  string
	output string
}

var kebabTests = []KebabTest{
	{"a", "a"},
	{"kebab", "kebab"},
	{"A", "a"},
	{"ID", "id"},
	{"MOTD", "motd"},
	{"Kebab", "kebab"},
	{"KebabTest", "kebab-test"},
	{"APIResponse", "api-response"},
	{"KebabID", "kebab-id"},
	{"Kebab_Id", "kebab-id"},
	{"Kebab_ID", "kebab-id"},
	{"KebabIDGoogle", "kebab-id-google"},
	{"LinuxMOTD", "linux-motd"},
	{"OMGWTFBBQ", "omgwtfbbq"},
	{"omg_wtf_bbq", "omg-wtf-bbq"},
	{"woof_woof", "woof-woof"},
	{"_woof_woof", "woof-woof"},
	{"woof_woof_", "woof-woof"},
	{"WOOF", "woof"},
	{"Woof", "woof"},
	{"woof", "woof"},
	{"woof0_woof1", "woof0-woof1"},
	{"_woof0_woof1_2", "woof0-woof1-2"},
	{"woof0_WOOF1_2", "woof0-woof1-2"},
	{"WOOF0", "woof0"},
	{"Woof1", "woof1"},
	{"woof2", "woof2"},
	{"woofWoof", "woof-woof"},
	{"woofWOOF", "woof-woof"},
	{"woof_WOOF", "woof-woof"},
	{"Woof_WOOF", "woof-woof"},
	{"WOOFWoofWoofWOOFWoofWoof", "woof-woof-woof-woof-woof-woof"},
	{"WOOF_Woof_woof_WOOF_Woof_woof", "woof-woof-woof-woof-woof-woof"},
	{"Woof_W", "woof-w"},
	{"Woof_w", "woof-w"},
	{"WoofW", "woof-w"},
	{"Woof_W_", "woof-w"},
	{"Woof_w_", "woof-w"},
	{"WoofW_", "woof-w"},
	{"WOOF_", "woof"},
	{"W_Woof", "w-woof"},
	{"w_Woof", "w-woof"},
	{"WWoof", "w-woof"},
	{"_W_Woof", "w-woof"},
	{"_w_Woof", "w-woof"},
	{"_WWoof", "w-woof"},
	{"_WOOF", "woof"},
	{"_woof", "woof"},
}

func TestKebabCase(t *testing.T) {
	for _, test := range kebabTests {
		if KebabCase(test.input) != test.output {
			t.Errorf("KebabCase(%q) -> %q, want %q", test.input, KebabCase(test.input), test.output)
		}
	}
}

func BenchmarkKebabCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range kebabTests {
			KebabCase(test.input)
		}
	}
}
